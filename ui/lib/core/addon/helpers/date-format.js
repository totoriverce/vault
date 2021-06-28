import { helper } from '@ember/component/helper';
import { format, parseISO } from 'date-fns';

export function dateFormat([date, style]) {
  // see format breaking in upgrade to date-fns 2.x https://github.com/date-fns/date-fns/blob/master/CHANGELOG.md#changed-5
  let number = typeof date === 'string' ? parseISO(date) : date;
  if (!number) {
    return;
  }
  if (number.toString().length === 10) {
    number = new Date(number * 1000);
  }
  return format(number, style);
}

export default helper(dateFormat);
