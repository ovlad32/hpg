package sparse

import "fmt"

/**
 *  Performs a logical <b>XOR</b> of the addressed target bit with the
 *  argument value. This bit set is modified so that the addressed bit has the
 *  value <code>true</code> if and only one of the following statements holds:
 *  <ul>
 *  <li>The addressed bit initially had the value <code>true</code>, and the
 *      value of the argument is <code>false</code>.
 *  <li>The bit initially had the value <code>false</code>, and the
 *      value of the argument is <code>true</code>.
 * </ul>
 *
 * @param       i a bit index
 * @param       value a boolean value to <b>XOR</b> with that bit
 * @exception   java.lang.IndexOutOfBoundsException if the specified index
 *              is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
//public void xor(int i, boolean value) {
func (this *BitSet) XorBit(i int32, value bool) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	if value {
		this.FlipBit(i)
	}
}

/**
 *  Performs a logical <b>XOR</b> of this bit set with the bit set argument
 *  within the given range. This resulting bit set is computed so that a bit
 *  within the range in it has the value <code>true</code> if and only if one
 *  of the following statements holds:
 *  <ul>
 *  <li>The bit initially had the value <code>true</code>, and the
 *      corresponding bit in the argument set has the value <code>false</code>.
 *  <li>The bit initially had the value <code>false</code>, and the
 *      corresponding bit in the argument set has the value <code>true</code>.
 * </ul>
 *  Outside the range this set is not changed.
 *
 * @param       i index of the first bit to be included in the operation
 * @param       j index after the last bit to included in the operation
 * @param       b the SparseBitSet with which to perform the <b>XOR</b>
 *              operation with this SparseBitSet
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative,
 *              or <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
//public void xor(int i, int j, SparseBitSet b) throws IndexOutOfBoundsException{
func (this *BitSet) XorRangeBitSet(i, j int32, b *BitSet) {
	this.setScanner(i, j, b, xorStrategy)
}

/**
 *  Performs a logical <b>XOR</b> of this bit set with the bit set argument.
 *  This resulting bit set is computed so that a bit in it has the value
 *  <code>true</code> if and only if one of the following statements holds:
 *  <ul>
 *  <li>The bit initially had the value <code>true</code>, and the
 *      corresponding bit in the argument set has the value <code>false</code>.
 *  <li>The bit initially had the value <code>false</code>, and the
 *      corresponding bit in the argument set has the value <code>true</code>.
 * </ul>
 *
 * @param       b the SparseBitSet with which to perform the <b>XOR</b>
 *              operation with thisSparseBitSet
 * @since       1.6
 */
//public void xor(SparseBitSet b) {

func (this *BitSet) XorBitSet(b *BitSet) {
	this.setScanner(0, b.bitsLength, b, xorStrategy)
}

/**
 * Performs a logical <b>XOR</b> of the two given <code>SparseBitSet</code>s.
 *  The resulting bit set is created so that a bit in it has the value
 *  <code>true</code> if and only if one of the following statements holds:
 *  <ul>
 *  <li>A bit in the first argument has the value <code>true</code>, and the
 *      corresponding bit in the second argument has the value
 *      <code>false</code>.</li>
 *  <li>A bit in the first argument has the value <code>false</code>, and the
 *      corresponding bit in the second argument has the value
 *      <code>true</code>.</li></ul>
 *
 * @param       a a SparseBitSet
 * @param       b another SparseBitSet
 * @return      a new SparseBitSet representing the <b>XOR</b> of the two sets
 * @since       1.6
 */
//public static SparseBitSet xor(SparseBitSet a, SparseBitSet b){
func Xor(a, b *BitSet) *BitSet {
	result := a.clone()
	result.XorBitSet(b)
	return result
}
