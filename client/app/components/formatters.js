import {makeSizeFormatter} from 'common-utils/formatters';

export const disk_mb = makeSizeFormatter('MiB', {
  precision:          1,
  precision_hundreds: 0,
});
export const disk_gb = makeSizeFormatter('GB', { //currently compatible with volume's formatters, but shouldn't we use GiB, precision = 1 instead?
  precision: 0,
});
export const memory_b = makeSizeFormatter('B', {
  precision:        1,
  use_binary_units: true,
});
export const memory_mb = makeSizeFormatter('MB', {
  precision:        1,
  use_binary_units: true,
});
export const memory_mib = makeSizeFormatter('MiB', {
  precision:        1,
  use_binary_units: true,
});
export const rate_mbs = makeSizeFormatter('Mb/s');
