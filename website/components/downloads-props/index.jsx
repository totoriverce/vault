import Link from 'next/link'
import Button from '@hashicorp/react-button'
import { VERSION, CHANGELOG_URL } from 'data/version'
import s from '../../pages/downloads/style.module.css'

const changelogUrl = CHANGELOG_URL.length
  ? CHANGELOG_URL
  : `https://github.com/hashicorp/vault/blob/v${VERSION}/CHANGELOG.md`

export default function DownloadsProps(preMerchandisingSlot) {
  return {
    changelog: changelogUrl,
    getStartedDescription:
      'Follow step-by-step tutorials on the essentials of Vault.',
    getStartedLinks: [
      {
        label: 'Getting Started with the CLI',
        href: 'http://learn.hashicorp.com/collections/vault/getting-started',
      },
      {
        label: 'Getting Started with Vault UI',
        href: 'http://learn.hashicorp.com/collections/vault/getting-started-ui',
      },
      {
        label: 'Vault on HCP',
        href: 'http://learn.hashicorp.com/collections/vault/getting-started-ui',
      },
      {
        label: 'View all Vault tutorials',
        href: 'https://learn.hashicorp.com/vault',
      },
    ],
    logo: (
      <img
        className={s.logo}
        alt="Vault"
        src={require('./img/vault-logo.svg')}
      />
    ),
    tutorialLink: {
      href: 'https://learn.hashicorp.com/vault',
      label: 'View Tutorials at HashiCorp Learn',
    },
    merchandisingSlot: (
      <>
        {preMerchandisingSlot && preMerchandisingSlot}

        <div className={s.merchandisingSlot}>
          <div className={s.centerWrapper}>
            <p>
              Want all of the power and security of Vault, without the
              complexity and overhead of managing it yourself?
            </p>
            <Button
              title="Sign up for HCP Vault"
              linkType="inbound"
              url="https://portal.cloud.hashicorp.com/sign-up?utm_source=vault_io&utm_content=download_cta"
              theme={{
                variant: 'tertiary',
                brand: 'vault',
              }}
            />
          </div>
        </div>

        <p className={s.releaseNote}>
          Release notes are available in our{' '}
          <Link href={`/docs/release-notes`}>
            <a>documentation</a>
          </Link>
          .
        </p>
      </>
    ),
  }
}
