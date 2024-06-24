export function one() {
    return two();
}

function two() {
    return "module one";
}
