package sparse

import "fmt"

/**
 *  Sets the bit at the specified index to <code>false</code>.
 *
 * @param       i a bit index.
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE.
 * @since       1.6
 */
func (bs *BitSet) Clear(i int32) {
	/*  In the interests of speed, no check is made here on whether the
	level3 block goes to all zero. This may be found and corrected
	in some later operation. */
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	if i >= bs.bitsLength {
		return
	}
	w := i >> cShift3
	a2 := bs.bits[w>>cShift1]
	if a2 == nil {
		return
	}
	a3 := a2[(w>>cShift2)&cMask2]
	if a3 == nil {
		return
	}
	a3[(w & cMask3)] &= ^wordType(uint(1) << remainderOf64(i)) //  Clear the indicated bit
	bs.cache.hash = 0                                          //  Invalidate size, etc.,
}

/**
 *  Sets the bits from the specified <code>i</code> (inclusive) to the
 *  specified <code>j</code> (exclusive) to <code>false</code>.
 *
 * @param       i index of the first bit to be cleared
 * @param       j index after the last bit to be cleared
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative,
 *              or <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
func (bs *BitSet) ClearRange(i, j int32) {
	bs.setScanner(i, j, nil, clearStrategy)
}

/**
 *  Sets all of the bits in this <code>SparseBitSet</code> to
 *  <code>false</code>.
 *
 * @since       1.6
 */
