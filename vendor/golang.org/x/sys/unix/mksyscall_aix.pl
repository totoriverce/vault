#!/usr/bin/env perl
# Copyright 2018 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This program reads a file containing function prototypes
# (like syscall_aix.go) and generates system call bodies.
# The prototypes are marked by lines beginning with "//sys"
# and read like func declarations if //sys is replaced by func, but:
#	* The parameter lists must give a name for each argument.
#	  This includes return parameters.
#	* The parameter lists must give a type for each argument:
#	  the (x, y, z int) shorthand is not allowed.
#	* If the return parameter is an error number, it must be named err.
#	* If go func name needs to be different than its libc name,
#	* or the function is not in libc, name could be specified
#	* at the end, after "=" sign, like
#	  //sys getsockopt(s int, level int, name int, val uintptr, vallen *_Socklen) (err error) = libsocket.getsockopt

use strict;

my $cmdline = "mksyscall_aix.pl " . join(' ', @ARGV);
my $errors = 0;
my $_32bit = "";
my $tags = "";  # build tags
my $aix = 0;
my $solaris = 0;

binmode STDOUT;

if($ARGV[0] eq "-b32") {
	$_32bit = "big-endian";
	shift;
} elsif($ARGV[0] eq "-l32") {
	$_32bit = "little-endian";
	shift;
}
if($ARGV[0] eq "-aix") {
	$aix = 1;
	shift;
}
if($ARGV[0] eq "-tags") {
	shift;
	$tags = $ARGV[0];
	shift;
}

if($ARGV[0] =~ /^-/) {
	print STDERR "usage: mksyscall_aix.pl [-b32 | -l32] [-tags x,y] [file ...]\n";
	exit 1;
}

sub parseparamlist($) {
	my ($list) = @_;
	$list =~ s/^\s*//;
	$list =~ s/\s*$//;
	if($list eq "") {
		return ();
	}
	return split(/\s*,\s*/, $list);
}

sub parseparam($) {
	my ($p) = @_;
	if($p !~ /^(\S*) (\S*)$/) {
		print STDERR "$ARGV:$.: malformed parameter: $p\n";
		$errors = 1;
		return ("xx", "int");
	}
	return ($1, $2);
}

my $package = "";
my $text = "";
my $c_extern = "/*\n#include <stdint.h>\n";
my @vars = ();
while(<>) {
	chomp;
	s/\s+/ /g;
	s/^\s+//;
	s/\s+$//;
	$package = $1 if !$package && /^package (\S+)$/;
	my $nonblock = /^\/\/sysnb /;
	next if !/^\/\/sys / && !$nonblock;

	# Line must be of the form
	# func Open(path string, mode int, perm int) (fd int, err error)
	# Split into name, in params, out params.
	if(!/^\/\/sys(nb)? (\w+)\(([^()]*)\)\s*(?:\(([^()]+)\))?\s*(?:=\s*(?:(\w*)\.)?(\w*))?$/) {
		print STDERR "$ARGV:$.: malformed //sys declaration\n";
		$errors = 1;
		next;
	}
	my ($nb, $func, $in, $out, $modname, $sysname) = ($1, $2, $3, $4, $5, $6);

	# Split argument lists on comma.
	my @in = parseparamlist($in);
	my @out = parseparamlist($out);

	$in = join(', ', @in);
	$out = join(', ', @out);

	# Try in vain to keep people from editing this file.
	# The theory is that they jump into the middle of the file
	# without reading the header.
	$text .= "// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT\n\n";

	# Check if value return, err return available
	my $errvar = "";
	my $retvar = "";
	my $rettype = "";
	foreach my $p (@out) {
		my ($name, $type) = parseparam($p);
		if($type eq "error") {
			$errvar = $name;
		} else {
			$retvar = $name;
			$rettype = $type;
		}
	}

	# System call name.
	#if($func ne "fcntl") {

	if($sysname eq "") {
		$sysname = "$func";
	}

	$sysname =~ s/([a-z])([A-Z])/${1}_$2/g;
	$sysname =~ y/A-Z/a-z/; # All libc functions are lowercase.

	my $C_rettype = "";
	if($rettype eq "unsafe.Pointer") {
		$C_rettype = "uintptr_t";
	} elsif($rettype eq "uintptr") {
		$C_rettype = "uintptr_t";
	} elsif($rettype =~ /^_/) {
		$C_rettype = "uintptr_t";
	} elsif($rettype eq "int") {
		$C_rettype = "int";
	} elsif($rettype eq "int32") {
		$C_rettype = "int";
	} elsif($rettype eq "int64") {
		$C_rettype = "long long";
	} elsif($rettype eq "uint32") {
		$C_rettype = "unsigned int";
	} elsif($rettype eq "uint64") {
		$C_rettype = "unsigned long long";
	} else {
		$C_rettype = "int";
	}
	if($sysname eq "exit") {
		$C_rettype = "void";
	}

	# Change types to c
	my @c_in = ();
	foreach my $p (@in) {
		my ($name, $type) = parseparam($p);
		if($type =~ /^\*/) {
			push @c_in, "uintptr_t";
			} elsif($type eq "string") {
			push @c_in, "uintptr_t";
		} elsif($type =~ /^\[\](.*)/) {
			push @c_in, "uintptr_t", "size_t";
		} elsif($type eq "unsafe.Pointer") {
			push @c_in, "uintptr_t";
		} elsif($type eq "uintptr") {
			push @c_in, "uintptr_t";
		} elsif($type =~ /^_/) {
			push @c_in, "uintptr_t";
		} elsif($type eq "int") {
			push @c_in, "int";
		} elsif($type eq "int32") {
			push @c_in, "int";
		} elsif($type eq "int64") {
			push @c_in, "long long";
		} elsif($type eq "uint32") {
			push @c_in, "unsigned int";
		} elsif($type eq "uint64") {
			push @c_in, "unsigned long long";
		} else {
			push @c_in, "int";
		}
	}

	if ($func ne "fcntl" && $func ne "FcntlInt" && $func ne "readlen" && $func ne "writelen") {
		# Imports of system calls from libc
		$c_extern .= "$C_rettype $sysname";
		my $c_in = join(', ', @c_in);
		$c_extern .= "($c_in);\n";
	}

	# So file name.
	if($aix) {
		if($modname eq "") {
			$modname = "libc.a/shr_64.o";
		} else {
			print STDERR "$func: only syscall using libc are available\n";
			$errors = 1;
			next;
		}
	}

	my $strconvfunc = "C.CString";
	my $strconvtype = "*byte";

	# Go function header.
	if($out ne "") {
		$out = " ($out)";
	}
	if($text ne "") {
		$text .= "\n"
	}

	$text .= sprintf "func %s(%s)%s {\n", $func, join(', ', @in), $out ;

	# Prepare arguments to call.
	my @args = ();
	my $n = 0;
	my $arg_n = 0;
	foreach my $p (@in) {
		my ($name, $type) = parseparam($p);
		if($type =~ /^\*/) {
			push @args, "C.uintptr_t(uintptr(unsafe.Pointer($name)))";
		} elsif($type eq "string" && $errvar ne "") {
			$text .= "\t_p$n := uintptr(unsafe.Pointer($strconvfunc($name)))\n";
			push @args, "C.uintptr_t(_p$n)";
			$n++;
		} elsif($type eq "string") {
			print STDERR "$ARGV:$.: $func uses string arguments, but has no error return\n";
			$text .= "\t_p$n := uintptr(unsafe.Pointer($strconvfunc($name)))\n";
			push @args, "C.uintptr_t(_p$n)";
			$n++;
		} elsif($type =~ /^\[\](.*)/) {
			# Convert slice into pointer, length.
			# Have to be careful not to take address of &a[0] if len == 0:
			# pass nil in that case.
			$text .= "\tvar _p$n *$1\n";
			$text .= "\tif len($name) > 0 {\n\t\t_p$n = \&$name\[0]\n\t}\n";
			push @args, "C.uintptr_t(uintptr(unsafe.Pointer(_p$n)))";
			$n++;
			$text .= "\tvar _p$n int\n";
			$text .= "\t_p$n = len($name)\n";
			push @args, "C.size_t(_p$n)";
			$n++;
		} elsif($type eq "int64" && $_32bit ne "") {
			if($_32bit eq "big-endian") {
				push @args, "uintptr($name >> 32)", "uintptr($name)";
			} else {
				push @args, "uintptr($name)", "uintptr($name >> 32)";
			}
			$n++;
		} elsif($type eq "bool") {
			$text .= "\tvar _p$n uint32\n";
			$text .= "\tif $name {\n\t\t_p$n = 1\n\t} else {\n\t\t_p$n = 0\n\t}\n";
			push @args, "_p$n";
			$n++;
		} elsif($type =~ /^_/) {
			push @args, "C.uintptr_t(uintptr($name))";
		} elsif($type eq "unsafe.Pointer") {
			push @args, "C.uintptr_t(uintptr($name))";
		} elsif($type eq "int") {
			if (($arg_n == 2) && (($func eq "readlen") || ($func eq "writelen"))) {
				push @args, "C.size_t($name)";
			} elsif ($arg_n == 0 && $func eq "fcntl") {
				push @args, "C.uintptr_t($name)";
			} elsif (($arg_n == 2) && (($func eq "fcntl") || ($func eq "FcntlInt"))) {
				push @args, "C.uintptr_t($name)";
			} else {
				push @args, "C.int($name)";
			}
		} elsif($type eq "int32") {
			push @args, "C.int($name)";
		} elsif($type eq "int64") {
			push @args, "C.longlong($name)";
		} elsif($type eq "uint32") {
			push @args, "C.uint($name)";
		} elsif($type eq "uint64") {
			push @args, "C.ulonglong($name)";
		} elsif($type eq "uintptr") {
			push @args, "C.uintptr_t($name)";
		} else {
			push @args, "C.int($name)";
		}
		$arg_n++;
	}
	my $nargs = @args;


	# Determine which form to use; pad args with zeros.
	if ($nonblock) {
	}

	my $args = join(', ', @args);
	my $call = "";
	if ($sysname eq "exit") {
		if ($errvar ne "") {
			$call .= "er :=";
		} else {
			$call .= "";
		}
	}  elsif ($errvar ne "") {
		$call .= "r0,er :=";
	}  elsif ($retvar ne "") {
		$call .= "r0,_ :=";
	}  else {
		$call .= ""
	}
	$call .= "C.$sysname($args)";

	# Assign return values.
	my $body = "";
	my $failexpr = "";

	for(my $i=0; $i<@out; $i++) {
		my $p = $out[$i];
		my ($name, $type) = parseparam($p);
		my $reg = "";
		if($name eq "err") {
			$reg = "e1";
		} else {
			$reg = "r0";
		}
		if($reg ne "e1" ) {
						$body .= "\t$name = $type($reg)\n";
		}
	}

	# verify return
	if ($sysname ne "exit" && $errvar ne "") {
		if ($C_rettype =~ /^uintptr/) {
			$body .= "\tif \(uintptr\(r0\) ==\^uintptr\(0\) && er != nil\) {\n";
			$body .= "\t\t$errvar = er\n";
			$body .= "\t}\n";
		} else {
			$body .= "\tif \(r0 ==-1 && er != nil\) {\n";
			$body .= "\t\t$errvar = er\n";
			$body .= "\t}\n";
		}
	} elsif ($errvar ne "") {
		$body .= "\tif \(er != nil\) {\n";
		$body .= "\t\t$errvar = er\n";
		$body .= "\t}\n";
	}

	$text .= "\t$call\n";
	$text .= $body;

	$text .= "\treturn\n";
	$text .= "}\n";
}

if($errors) {
	exit 1;
}

print <<EOF;
// $cmdline
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build $tags

package $package


$c_extern
*/
import "C"
import (
	"unsafe"
	"syscall"
)


EOF

print "import \"golang.org/x/sys/unix\"\n" if $package ne "unix";

chomp($_=<<EOF);

$text
EOF
print $_;
exit 0;
