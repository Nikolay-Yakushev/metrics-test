package helpers


func Contains[T comparable](elems []T, val T) bool {
    for ix := range elems {
        if val == elems[ix] {
            return true
        }
    }
    return false
}
