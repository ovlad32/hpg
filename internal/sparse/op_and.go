package sparse

import "fmt"

/**
 *  Performs a logical <b>AND</b> of the addressed target bit with the argument
 *  value. This bit set is modified so that the addressed bit has the value
 *  <code>true</code> if and only if it both initially had the value
 *  <code>true</code> and the argument value is also <code>true</code>.
 *
 * @param       i a bit index
 * @param       value a boolean value to <b>AND</b> with that bit
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
//public void and(int i, boolean value) throws IndexOutOfBoundsException
func (this *BitSet) AndBit(i int32, value bool) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	if !value {
		this.Clear(i)
	}
}

/**
 *  Performs a logical <b>AND</b> of this target bit set with the argument bit
 *  set within the given range of bits. Within the range, this bit set is
 *  modified so that each bit in it has the value <code>true</code> if and only
 *  if it both initially had the value <code>true</code> and the corresponding
 *  bit in the bit set argument also had the value <code>true</code>. Outside
 *  the range, this set is not changed.
 *
 * @param       i index of the first bit to be included in the operation
 * @param       j index after the last bit to included in the operation
 * @param       b a SparseBitSet
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative,
 *              or <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
//public void and(int i, int j, SparseBitSet b) throws IndexOutOfBoundsException
func (this *BitSet) AndRangeBitSet(i, j int32, b *BitSet) {
	this.setScanner(i, j, b, andStrategyType{})
}

/**
 *  Performs a logical <b>AND</b> of this target bit set with the argument bit
 *  set. This bit set is modified so that each bit in it has the value
 *  <code>true</code> if and only if it both initially had the value
 *  <code>true</code> and the corresponding bit in the bit set argument also
 *  had the value <code>true</code>.
 *
 * @param       b a SparseBitSet
 * @since       1.6
 */
//    public void and(SparseBitSet b)
func (this *BitSet) AndBitSet(b *BitSet) {
	{
		bmin := len(this.bits)
		if len(b.bits) < bmin {
			bmin = len(b.bits)
		}
		this.nullify(int32(bmin)) // Optimisation
	}
	{
		bmin := this.bitsLength
		if b.bitsLength < bmin {
			bmin = b.bitsLength
		}
		this.setScanner(0, bmin, b, andStrategy)
	}
}

/**
 *  Performs a logical <b>AND</b> of the two given <code>SparseBitSet</code>s.
 *  The returned <code>SparseBitSet</code> is created so that each bit in it
 *  has the value <code>true</code> if and only if both the given sets
 *  initially had the corresponding bits <code>true</code>, otherwise
 *  <code>false</code>.
 *
 * @param       a a SparseBitSet
 * @param       b another SparseBitSet
 * @return      a new SparseBitSet representing the <b>AND</b> of the two sets
 * @since       1.6
 */
//public static SparseBitSet and(SparseBitSet a, SparseBitSet b)
func And(a, b *BitSet) *BitSet {
	result := a.clone()
	result.AndBitSet(b)
	return result
}
