package stringutil

//Reverse returns the reverse of the string passed as `s`
func Reverse(s string) string {
    //convert string to a slice so we can manipulate it
    //since strings are immutable, this creates a copy of
    //the string so we won't be modifying the original
    chars := []byte(s)

    //starting from each end, swap the values in the slice
    //elements, stopping when we get to the middle
    for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
        chars[i], chars[j] = chars[j], chars[i]
    }

    //return the reversed slice as a string
    return string(chars)
}