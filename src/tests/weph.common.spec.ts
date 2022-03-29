export function q4(lon: number): boolean {
    return lon >= 270.0 && lon < 360.0;
}

export function q1(lon: number): boolean {
    return lon >= 0.0 && lon < 90.0;
}
