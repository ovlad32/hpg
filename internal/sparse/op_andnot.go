package sparse

import "fmt"

/**
 *  Performs a logical <b>AndNOT</b> of the addressed target bit with the
 *  argument value. This bit set is modified so that the addressed bit has the
 *  value <code>true</code> if and only if it both initially had the value
 *  <code>true</code> and the argument value is <code>false</code>.
 *
 * @param       i a bit index
 * @param       value a boolean value to AndNOT with that bit
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
//public void andNot(int i, boolean value)
func (bs *BitSet) AndNotBit(i int32, value bool) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	if value {
		bs.Clear(i)
	}
}

/**
 *  Performs a logical <b>AndNOT</b> of this target bit set with the argument
 *  bit set within the given range of bits. Within the range, this bit set is
 *  modified so that each bit in it has the value <code>true</code> if and only
 *  if it both initially had the value <code>true</code> and the corresponding
 *  bit in the bit set argument has the value <code>false</code>. Outside
 *  the range, this set is not changed.
 *
 * @param       i index of the first bit to be included in the operation
 * @param       j index after the last bit to included in the operation
 * @param       b the SparseBitSet with which to mask this SparseBitSet
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative,
 *              or <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
//    public void andNot(int i, int j, SparseBitSet b)
func (bs *BitSet) AndNotRangeBitSet(i, j int32, b *BitSet) {
	bs.setScanner(i, j, b, andNotStrategy)
}

/**
 *  Performs a logical <b>AndNOT</b> of this target bit set with the argument
 *  bit set. This bit set is modified so that each bit in it has the value
 *  <code>true</code> if and only if it both initially had the value
 *  <code>true</code> and the corresponding bit in the bit set argument has
 *  the value <code>false</code>.
 *
 * @param       b the SparseBitSet with which to mask this SparseBitSet
 * @since       1.6
 */
//    public void andNot(SparseBitSet b)
func (bs *BitSet) AndNotBitSet(b *BitSet) {
	bmin := bs.bitsLength
	if b.bitsLength < bmin {
		bmin = b.bitsLength
	}
	bs.setScanner(0, bmin, b, andNotStrategy)
}

/**
 *  Creates a bit set from thie first <code>SparseBitSet</code> whose
 *  corresponding bits are cleared by the set bits of the second
 *  <code>SparseBitSet</code>. The resulting bit set is created so that a bit
 *  in it has the value <code>true</code> if and only if the corresponding bit
 *  in the <code>SparseBitSet</code> of the first is set, and that same
 *  corresponding bit is not set in the <code>SparseBitSet</code> of the second
 *  argument.
 *
 * @param a     a SparseBitSet
 * @param b     another SparseBitSet
 * @return      a new SparseBitSet representing the <b>AndNOT</b> of the
 *              two sets
 * @since       1.6
 */
//public static SparseBitSet andNot(SparseBitSet a, SparseBitSet b){
func AndNot(a, b *BitSet) *BitSet {
	result := a.clone()
	result.AndNotBitSet(b)
	return result
}
