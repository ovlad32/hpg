package sparse

import "fmt"

/**
 *  Performs a logical <b>OR</b> of the addressed target bit with the
 *  argument value. This bit set is modified so that the addressed bit has the
 *  value <code>true</code> if and only if it both initially had the value
 *  <code>true</code> or the argument value is <code>true</code>.
 *
 * @param       i a bit index
 * @param       value a boolean value to OR with that bit
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
//public void or(int i, boolean value)
func (bs *BitSet) OrBit(i int32, value bool) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	if value {
		bs.Set(i)
	}
}

/**
 *  Performs a logical <b>OR</b> of the addressed target bit with the
 *  argument value within the given range. This bit set is modified so that
 *  within the range a bit in it has the value <code>true</code> if and only if
 *  it either already had the value <code>true</code> or the corresponding bit
 *  in the bit set argument has the value <code>true</code>. Outside the range
 *  this set is not changed.
 *
 * @param       i index of the first bit to be included in the operation
 * @param       j index after the last bit to included in the operation
 * @param       b the SparseBitSet with which to perform the <b>OR</b>
 *              operation with this SparseBitSet
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative,
 *              or <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
//public void or(int i, int j, SparseBitSet b) throws IndexOutOfBoundsException
func (bs *BitSet) OrRangeBitSet(i, j int32, b *BitSet) {
	bs.setScanner(i, j, b, orStrategy)
}

/**
 *  Performs a logical <b>OR</b> of this bit set with the bit set argument.
 *  This bit set is modified so that a bit in it has the value <code>true</code>
 *  if and only if it either already had the value <code>true</code> or the
 *  corresponding bit in the bit set argument has the value <code>true</code>.
 *
 * @param       b the SparseBitSet with which to perform the <b>OR</b>
 *              operation with this SparseBitSet
 * @since       1.6
 */
//public void or(SparseBitSet b){
func (bs *BitSet) OrBitSet(b *BitSet) {
	bs.setScanner(0, b.bitsLength, b, orStrategy)
}

/**
 *  Performs a logical <b>OR</b> of the two given <code>SparseBitSet</code>s.
 *  The returned <code>SparseBitSet</code> is created so that a bit in it has
 *  the value <code>true</code> if and only if it either had the value
 *  <code>true</code> in the set given by the first arguemetn or had the value
 *  <code>true</code> in the second argument, otherwise <code>false</code>.
 *
 * @param       a a SparseBitSet
 * @param       b another SparseBitSet
 * @return      new SparseBitSet representing the <b>OR</b> of the two sets
 * @since       1.6
 */
//public static SparseBitSet or(SparseBitSet a, SparseBitSet b) {
func Or(a, b *BitSet) *BitSet {
	result := a.clone()
	result.OrBitSet(b)
	return result
}
