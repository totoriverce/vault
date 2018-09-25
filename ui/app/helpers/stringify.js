import { helper as buildHelper } from '@ember/component/helper';

export function stringify([target], { skipFormat }) {
  if (skipFormat) {
    return JSON.stringify(target);
  }
  return JSON.stringify(target, null, 2);
}

export default buildHelper(stringify);
