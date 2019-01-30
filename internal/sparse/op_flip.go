package sparse

import "fmt"

/**
 *  Sets the bit at the specified index to the complement of its current value.
 *
 * @param       i the index of the bit to flip
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (bs *BitSet) FlipBit(i int32) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	w := i >> cShift3
	w1 := w >> cShift1
	w2 := (w >> cShift2) & cMask2

	if i >= bs.bitsLength {
		bs.resize(i)
	}

	var a2 b2DimType
	var a3 b1DimType
	if a2 = bs.bits[w1]; a2 == nil {
		a2 = make(b2DimType, cLength2)
		a3 = make(b1DimType, cLength3)
		a2[w2] = a3
	} else {
		if a3 = a2[w2]; a3 == nil {
			a3 = make(b1DimType, cLength3)
			a2[w2] = a3
		}
	}
	a3[(w & cMask3)] = a3[(w&cMask3)] ^ wordType(uint(1)<<remainderOf64(i)) //Flip the designated bit
	bs.cache.hash = 0                                                       //  Invalidate size, etc., values
}

/**
 *  Sets each bit from the specified <code>i</code> (inclusive) to the
 *  specified <code>j</code> (exclusive) to the complement of its current
 *  value.
 *
 * @param       i index of the first bit to flip
 * @param       j index after the last bit to flip
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or is
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative, or
 *              <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
func (bs *BitSet) FlipRange(i, j int32) {
	bs.setScanner(i, j, nil, flipStrategy)
}
