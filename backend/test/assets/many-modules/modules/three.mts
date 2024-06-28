import { two } from './two.mts';
import { common } from '../../common/common.mts';

export function three() {
    common();
    return "Module three called " + two();
}
