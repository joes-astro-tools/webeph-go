declare namespace jasmine {
    interface Matchers<T> {
        toBeWithinTolerance(expected: T, tolerance: T): boolean;
    }
}
