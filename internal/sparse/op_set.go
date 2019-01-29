import "fmt"

/**
 *  Sets the bit at the specified index.
 *
 * @param       i a bit index
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (this *BitSet) Set(i int32) {
	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	w := i >> SHIFT3
	w1 := w >> SHIFT1
	w2 := (w >> SHIFT2) & MASK2

	if i >= this.bitsLength {
		this.resize(i)
	}
	a2 := this.bits[w1]
	if a2 == nil {
		a2 = make(b2DimType, LENGTH2)
		this.bits[w1] = a2
	}

	a3 := a2[w2]
	if a3 == nil {
		a3 = make(b1DimType, LENGTH3)
		a2[w2] = a3
	}
	a3[(w & MASK3)] |= wordType(uint(1) << remainderOf64(i))
	this.cache.hash = 0 //Invalidate size, etc., scan
}

/**
 *  Sets the bit at the specified index to the specified value.
 *
 * @param       i a bit index
 * @param       value a boolean value to set
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (this *BitSet) SetBit(i int32, value bool) {
	if value {
		this.Set(i)
	} else {
		this.Clear(i)
	}
}

/**
 *  Sets the bits from the specified <code>i</code> (inclusive) to the specified
 *  <code>j</code> (exclusive) to <code>true</code>.
 *
 * @param       i index of the first bit to be set
 * @param       j index after the last bit to be se
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or is
 *              equal to Integer.MAX_INT, or <code>j</code> is negative, or
 *              <code>i</code> is larger than <code>j</code>.
 * @since       1.6
 */
func (this *BitSet) SetRange(i, j int32) {
	this.setScanner(i, j, nil, setStrategy)
}

/**
 *  Sets the bits from the specified <code>i</code> (inclusive) to the specified
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
func (this *BitSet) SetRangeBit(i, j int32, value bool) {
	if value {
		this.SetRange(i, j)
	} else {
		this.ClearRange(i, j)
	}
}