package sparse

import "fmt"

/*Set - sets the bit at the specified index.
 *
 *
 * @param       i a bit index
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (bs *BitSet) Set(i int32) {
	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	w := i >> cShift3
	w1 := w >> cShift1
	w2 := (w >> cShift2) & cMask2

	if i >= bs.bitsLength {
		bs.resize(i)
	}
	a2 := bs.bits[w1]
	if a2 == nil {
		a2 = make(b2DimType, cLength2)
		bs.bits[w1] = a2
	}

	a3 := a2[w2]
	if a3 == nil {
		a3 = make(b1DimType, cLength3)
		a2[w2] = a3
	}
	a3[(w & cMask3)] |= wordType(uint(1) << remainderOf64(i))
	bs.cache.hash = 0 //Invalidate size, etc., scan
}

/*SetBit - sets the bit at the specified index to the specified value.
 *
 * @param       i a bit index
 * @param       value a boolean value to set
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (bs *BitSet) SetBit(i int32, value bool) {
	if value {
		bs.Set(i)
	} else {
		bs.Clear(i)
	}
}

/*SetRange - sets the bits from the specified <code>i</code> (inclusive) to the specified
 *  <code>j</code> (exclusive) to <code>true</code>.
 *
 * @param       i index of the first bit to be set
 * @param       j index after the last bit to be se
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or is
 *              equal to Integer.MAX_INT, or <code>j</code> is negative, or
 *              <code>i</code> is larger than <code>j</code>.
 * @since       1.6
 */
func (bs *BitSet) SetRange(i, j int32) {
	bs.setScanner(i, j, nil, setStrategy)
}

/*SetRangeBit - sets the bits from the specified <code>i</code> (inclusive) to the specified
 *  <code>j</code> (exclusive) to the specified value.
 *
 * @param       i index of the first bit to be set
 * @param       j index after the last bit to be set
 * @param       value to which to set the selected bits
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or is
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative, or
 *              <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
func (bs *BitSet) SetRangeBit(i, j int32, value bool) {
	if value {
		bs.SetRange(i, j)
	} else {
		bs.ClearRange(i, j)
	}
}
