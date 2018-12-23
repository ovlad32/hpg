package sparse

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"strconv"
)

type wordType uint64
type b3DimType [][][]wordType
type b2DimType [][]wordType
type b1DimType []wordType

const compactionCountDefault int32 = 2

/**
 *  The number of bits in a long value.
 */
const LENGTH4 int32 = 64 //Long.SIZE
const INTEGER_SIZE = 32  //Integer.SIZE

/**
 *  The number of bits in a positive integer, and the size of permitted index
 *  of a bit in the bit set.
 */
const INDEX_SIZE uint32 = INTEGER_SIZE - 1 //Integer.SIZE - 1;

/**
 *  The label (index) of a bit in the bit set is essentially broken into
 *  4 "levels". Respectively (from the least significant end), level4, the
 *  address within word, the address within a level3 block, the address within
 *  a level2 area, and the level1 address of that area within the set.
 *
 *  LEVEL4 is the number of bits of the level4 address (number of bits need
 *  to address the bits in a long)
 */
const LEVEL4 uint32 = 6

/**
 *  LEVEL3 is the number of bits of the level3 address.
 *  Do not change!
 */
const LEVEL3 uint32 = 5

/**
 *  LEVEL2 is the number of bits of the level3 address.
 *  Do not change!
 */
const LEVEL2 uint32 = 5

/**
 *  LEVEL1 is the number of bits of the level1 address.
 */
const LEVEL1 = INDEX_SIZE - LEVEL2 - LEVEL3 - LEVEL4

/**
 *  MAX_LENGTH1 is the maximum number of entries in the level1 set array.
 */
const MAX_LENGTH1 int32 = 1 << uint(LEVEL1)

/**
 *  LENGTH2 is the number of entries in the any level2 area.
 */
const LENGTH2 int32 = 1 << uint(LEVEL2)

/**
 *  LENGTH3 is the number of entries in the any level3 block.
 */
const LENGTH3 int32 = 1 << uint(LEVEL3)

/**
 *  The shift to create the word index. (I.e., move it to the right end)
 */
const SHIFT3 uint32 = LEVEL4

/**
 *  MASK3 is the mask to extract the LEVEL3 address from a word index
 *  (after shifting by SHIFT3).
 */
const MASK3 int32 = int32(LENGTH3 - 1)

/**
 *  SHIFT2 is the shift to bring the level2 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const SHIFT2 uint32 = LEVEL3

/**
 *  UNIT is the greatest number of bits that can be held in one level1 entry.
 *  That is, bits per word by words per level3 block by blocks per level2 area.
 */
const UNIT int32 = LENGTH2 * LENGTH3 * LENGTH4

/**
 *  MASK2 is the mask to extract the LEVEL2 address from a word index
 *  (after shifting by SHIFT3 and SHIFT2).
 */
const MASK2 int32 = int32(LENGTH2 - 1)

/**
 *  SHIFT1 is the shift to bring the level1 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const SHIFT1 uint32 = LEVEL2 + LEVEL3

/**
 *  LENGTH2_SIZE is maximum index of a LEVEL2 page.
 */
const LENGTH2_SIZE int32 = int32(LENGTH2 - 1)

/**
 *  LENGTH3_SIZE is maximum index of a LEVEL3 page.
 */
const LENGTH3_SIZE int32 = int32(LENGTH3 - 1)

/**
 *  LENGTH4_SIZE is maximum index of a bit in a LEVEL4 word.
 */
const LENGTH4_SIZE int32 = int32(LENGTH4 - 1)

/** An empty level 3 block is kept for use when scanning. When a source block
 *  is needed, and there is not already one in the corresponding bit set, the
 *  ZERO_BLOCK is used (as a read-only block). It is a source of zero values
 *  so that code does not have to test for a null level3 block. This is a
 *  static block shared everywhere.
 */
var ZERO_BLOCK = make(b1DimType, LENGTH3)

/**
 *  Word and block <b>and</b> strategy.
 */
var andStrategy = new(andStrategyType)

/**
 *  Word and block <b>andNot</b> strategy.
 */
var andNotStrategy = new(andNotStrategyType)

/**
 *  Word and block <b>clear</b> strategy.
 */
var clearStrategy = new(clearStrategyType)

/**
 *  Word and block <b>copy</b> strategy.
 */
var copyStrategy = new(copyStrategyType)

/**
 *  Word and block <b>flip</b> strategy.
 */
var flipStrategy = new(flipStrategyType)

/**
 *  Word and block <b>intersects</b> strategy.
 */
var intersectsStrategy = new(intersectsStrategyType)

/**
 *  Word and block <b>or</b> strategy.
 */
var orStrategy = new(orStrategyType)

/**
 *  Word and block <b>set</b> strategy.
 */
var setStrategy = new(setStrategyType)

/**
 *  Word and block <b>xor</b> strategy.
 */
var xorStrategy = new(xorStrategyType)

type BitSet struct {
	/**
	 *  This value controls for format of the toString() output.
	 * @see #toStringCompaction(int)
	 */
	compactionCount int32

	/**
	 *  The storage for this SparseBitSet. The <i>i</i>th bit is stored in a word
	 *  represented by a long value, and is at bit position <code>i % 64</code>
	 *  within that word (where bit position 0 refers to the least significant bit
	 *  and 63 refers to the most significant bit).
	 *  <p>
	 *  The words are organized into blocks, and the blocks are accessed by two
	 *  additional levels of array indexing.
	 */
	bits b3DimType

	/**
	 *  For the current size of the bits array, this is the maximum possible
	 *  length of the bit set, i.e., the index of the last possible bit, plus one.
	 *  Note: this not the value returned by <i>length</i>().
	 * @see #resize(int)
	 * @see #length()
	 */
	bitsLength int32
	/**
	 *  Holds reference to the cache of statistics values computed by the
	 *  UpdateStrategy
	 * @see SparseBitSet.Cache
	 * @see SparseBitSet.UpdateStrategy
	 */
	cache *cacheType
	/**
	 *  A spare level 3 block is kept for use when scanning. When a target block
	 *  is needed, and there is not already one in the bit set, the spare is
	 *  provided. If non-zero values are placed into this block, it is moved to the
	 *  resulting set, and a new spare is acquired. Note: a new spare needs to
	 *  be allocated when the set is cloned (so that the spare is not shared
	 *  between two sets).
	 */
	spare b1DimType

	/**
	 *  Word and block <b>equals</b> strategy.
	 */
	equalsStrategy *equalsStrategyType
	/**
	 *  Word and block <b>update</b> strategy.
	 */
	updateStrategy *updateStrategyType
}

/**
 *  Constructs an empty bit set with the default initial size.
 *  Initially all bits are effectively <code>false</code>.
 *
 * @since       1.6
 */
//    public SparseBitSet()
func New() *BitSet {
	return newWithSizeAndCompactionCount(1, compactionCountDefault)
}

/**
*  Creates a bit set whose initial size is large enough to efficiently
*  represent bits with indices in the range <code>0</code> through
*  at least <code>nbits-1</code>. Initially all bits are effectively
*  <code>false</code>.
*  <p>
*  No guarantees are given for how large or small the actual object will be.
*  The setting of bits above the given range is permitted (and will perhaps
*  eventually cause resizing).
*
* @param       nbits the initial provisional length of the SparseBitSet
* @throws      java.lang.NegativeArraySizeException if the specified initial
*              length is negative
* @see         #SparseBitSet()
* @since       1.
 */
// public SparseBitSet(int nbits) throws NegativeArraySizeException
func NewWithSize(capacity uint32) *BitSet {
	return newWithSizeAndCompactionCount(1, compactionCountDefault)
}

/**
 *  Constructor for a new (sparse) bit set. All bits initially are effectively
 *  <code>false</code>. This is a internal constructor that collects all the
 *  needed actions to initialise the bit set.
 *  <p>
 *  The capacity is taken to be a <i>suggestion</i> for a size of the bit set,
 *  in bits. An appropiate table size (a power of two) is then determined and
 *  used. The size will be grown as needed to accomodate any bits addressed
 *  during the use of the bit set.
 *
 * @param       capacity a size in terms of bits
 * @param       compactionCount the compactionCount to be inherited (for
 *              internal generation)
 * @exception   NegativeArraySizeException if the specified initial size
 *              is negative
 * @since       1.6
 */
//    protected SparseBitSet(int capacity, int compactionCount) throws NegativeArraySizeException
func newWithSizeAndCompactionCount(capacity int32, compactionCount int32) *BitSet {
	result := &BitSet{
		compactionCount: compactionCount,
	}
	result.resize(capacity - 1) //  Resize takes last usable index
	result.compactionCount = compactionCount
	/*  Ensure there is a spare level 3 block for the use of the set scanner.*/
	result.constructorHelper()
	result.statisticsUpdate()
	return result
}

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
	this.setScanner(i, j, b, andStrategy)
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
func (this *BitSet) AndNotBit(i int32, value bool) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	if value {
		this.Clear(i)
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
func (this *BitSet) AndNotRangeBitSet(i, j int32, b *BitSet) {
	this.setScanner(i, j, b, andNotStrategy)
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
func (this *BitSet) AndNotBitSet(b *BitSet) {
	bmin := this.bitsLength
	if b.bitsLength < bmin {
		bmin = b.bitsLength
	}
	this.setScanner(0, bmin, b, andNotStrategy)
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

/**
 *  Sets the bit at the specified index to the complement of its current value.
 *
 * @param       i the index of the bit to flip
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (this *BitSet) FlipBit(i int32) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	w := i >> SHIFT3
	w1 := w >> SHIFT1
	w2 := (w >> SHIFT2) & MASK2

	if i >= this.bitsLength {
		this.resize(i)
	}

	var a2 b2DimType
	var a3 b1DimType
	if a2 = this.bits[w1]; a2 == nil {
		a2 = make(b2DimType, LENGTH2)
		a3 = make(b1DimType, LENGTH3)
		a2[w2] = a3
	} else {
		if a3 = a2[w2]; a3 == nil {
			a3 = make(b1DimType, LENGTH3)
			a2[w2] = a3
		}
	}
	a3[(w & MASK3)] = a3[(w&MASK3)] ^ wordType(bits.RotateLeft(uint(1), int(i))) //Flip the designated bit
	this.cache.hash = 0                                                          //  Invalidate size, etc., values
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
func (this *BitSet) FlipRange(i, j int32) {
	this.setScanner(i, j, nil, flipStrategy)
}

/**
 *  Returns the value of the bit with the specified index. The value is
 *  <code>true</code> if the bit with the index <code>i</code> is currently set
 *  in this <code>SparseBitSet</code>; otherwise, the result is
 *  <code>false</code>.
 *
 * @param       i the bit index
 * @return      the boolean value of the bit with the specified index.
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */

/**
 *  Returns the value of the bit with the specified index. The value is
 *  <code>true</code> if the bit with the index <code>i</code> is currently set
 *  in this <code>SparseBitSet</code>; otherwise, the result is
 *  <code>false</code>.
 *
 * @param       i the bit index
 * @return      the boolean value of the bit with the specified index.
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (this *BitSet) GetBit(i int32) bool {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}

	w := i >> SHIFT3
	if i < this.bitsLength {
		a2 := this.bits[w>>SHIFT1]
		if a2 != nil {
			a3 := a2[(w>>SHIFT2)&MASK2]
			if a3 != nil {
				result := a3[(w&MASK3)] & (wordType(bits.RotateLeft(uint(1), int(i))))
				return result != 0
			}
		}
	}
	return false
}

/**
 *  Returns a new <code>SparseBitSet</code> composed of bits from this
 *  <code>SparseBitSet</code> from <code>i</code> (inclusive) to <code>j</code>
 *  (exclusive).
 *
 * @param       i index of the first bit to include
 * @param       j index after the last bit to include
 * @return      a new SparseBitSet from a range of this SparseBitSet
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or is
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative, or
 *              <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */

func (this *BitSet) GetBitSetFromRange(i, j int32) *BitSet {
	result := newWithSizeAndCompactionCount(j, this.compactionCount)
	result.setScanner(i, j, this, copyStrategy)
	return result
}

/**
 *  Returns a hash code value for this bit set. The hash code depends only on
 *  which bits have been set within this <code>SparseBitSet</code>. The
 *  algorithm used to compute it may be described as follows.
 *  <p>
 *  Suppose the bits in the <code>SparseBitSet</code> were to be stored in an
 *  array of <code>long</code> integers called, say, <code>bits</code>, in such
 *  a manner that bit <code>i</code> is set in the <code>SparseBitSet</code>
 *  (for nonnegative values of  <code>i</code>) if and only if the expression
 *  <pre>
 *  ((i&gt;&gt;6) &lt; bits.length) &amp;&amp; ((bits[i&gt;&gt;6] &amp; (1L &lt;&lt; (bit &amp; 0x3F))) != 0)
 *  </pre>
 *  is true. Then the following definition of the <code>hashCode</code> method
 *  would be a correct implementation of the actual algorithm:
 *  <pre>
 *  public int hashCode()
 *  {
 *      long hash = 1234L;
 *      for( int i = bits.length; --i &gt;= 0; )
 *          hash ^= bits[i] * (i + 1);
 *      return (int)((h &gt;&gt; 32) ^ h);
 *  }</pre>
 *  Note that the hash code values change if the set of bits is altered.
 *
 * @return      a hash code value for this bit set
 * @since       1.6
 * @see         Object#equals(Object)
 * @see         java.util.Hashtable
 */
/*@Override
 func(this *BitSet)  int hashCode()
 {
	 statisticsUpdate();
	 return cache.hash;
 }*/

/**
 *  Returns true if the specified <code>SparseBitSet</code> has any bits
 *  within the given range <code>i</code> (inclusive) to <code>j</code>
 *  (exclusive) set to <code>true</code> that are also set to <code>true</code>
 *  in the same range of this <code>SparseBitSet</code>.
 *
 * @param       i index of the first bit to include
 * @param       j index after the last bit to include
 * @param       b the SparseBitSet with which to intersect
 * @return      the boolean indicating whether this SparseBitSet intersects the
 *              specified SparseBitSet
 * @exception   IndexOutOfBoundsException if <code>i</code> is negative or
 *              equal to Integer.MAX_VALUE, or <code>j</code> is negative,
 *              or <code>i</code> is larger than <code>j</code>
 * @since       1.6
 */
//public boolean intersects(int i, int j, SparseBitSet b) throws IndexOutOfBoundsException {
/*func (this *BitSet) IntersectsRangeBitSet(i, j int32, b *BitSet) {
	this.setScanner(i, j, b, intersectsStrategy)
	return intersectsStrategy.result
}*/

/**
 *  Returns true if the specified <code>SparseBitSet</code> has any bits set to
 *  <code>true</code> that are also set to <code>true</code> in this
 *  <code>SparseBitSet</code>.
 *
 * @param       b a SparseBitSet with which to intersect
 * @return      boolean indicating whether this SparseBitSet intersects the
 *              specified SparseBitSet
 * @since       1.6
 */
/*public boolean intersects(SparseBitSet b)
 {
	 setScanner(0, Math.max(bitsLength, b.bitsLength), b, intersectsStrategy);
	 return intersectsStrategy.result;
 } */

/**
 *  Returns true if this <code>SparseBitSet</code> contains no bits that are
 *  set to <code>true</code>.
 *
 * @return      the boolean indicating whether this SparseBitSet is empty
 * @since       1.6
 */
func (this *BitSet) isEmpty() bool {
	this.statisticsUpdate()
	return this.cache.cardinality == 0
}

/**
 *  Returns the "logical length" of this <code>SparseBitSet</code>: the index
 *  of the highest set bit in the <code>SparseBitSet</code> plus one. Returns
 *  zero if the <code>SparseBitSet</code> contains no set bits.
 *
 * @return      the logical length of this SparseBitSet
 * @since       1.6
 */
func (this *BitSet) Length() int32 {
	this.statisticsUpdate()
	return this.cache.length
}

/**
 *  Returns the number of bits of space nominally in use by this
 *  <code>SparseBitSet</code> to represent bit values. The count of bits in
 *  the set is the (label of the last set bit) + 1 - (the label of the first
 *  set bit).
 *
 * @return      the number of bits (true and false) nominally in this bit set
 *              at this moment
 * @since       1.6
 */
func (this *BitSet) Size() int32 {
	this.statisticsUpdate()
	return this.cache.size
}

/**
 *  Returns the index of the first bit that is set to <code>false</code> that
 *  occurs on or after the specified starting index.
 *
 * @param       i the index to start checking from (inclusive)
 * @return      the index of the next clear bit, or -1 if there is no such bit
 * @exception   IndexOutOfBoundsException if the specified index is negative
 * @since       1.6
 */
func (this *BitSet) NextClearBit(i int32) int32 {
	/*  The index of this method is permitted to be Integer.MAX_VALUE, as this
	is needed to make this method work together with the method
	nextSetBit()--as might happen if a search for the next clear bit is
	started after finding a set bit labelled Integer.MAX_VALUE-1. This
	case is not optimised, the code will eventually return -1 (since
	the Integer.MAX_VALUEth bit does "exist," and is 0. */

	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}

	/*  This is the word from which the search begins. */
	w := i >> SHIFT3
	w3 := w & MASK3
	w2 := (w >> SHIFT2) & MASK2
	w1 := w >> SHIFT1

	nword := wordType(^int64(0) << uint(i))
	aLength := int32(len(this.bits))

	/*  Is the next clear bit in the same word at the nominated beginning bit
	(including the nominated beginning bit itself). The first check is
	whether the starting bit is within the structure at all. */

	if w1 < aLength {
		if a2 := this.bits[w1]; a2 != nil {
			if a3 := a2[w2]; a3 != nil {
				if nword = ^a3[w3] & wordType(int64(0)<<uint(i)); nword == 0 {
					w++
					w3 = w & MASK3
					w2 = (w >> SHIFT2) & MASK2
					w1 = w >> SHIFT1
					nword = ^wordType(0)
				loop:
					for ; w1 != aLength; w1++ {
						if a2 = this.bits[w1]; a2 == nil {
							break
						}
						for ; w2 != LENGTH2; w2++ {
							if a3 = a2[w2]; a3 == nil {
								break loop
							}
							for ; w3 != LENGTH3; w3++ {
								if nword = ^a3[w3]; nword != 0 {
									break loop
								}
							}
							w3 = 0
						}
						w2, w3 = 0, 0
					}
				}
			}
		}

		/*  So now start a search though the rest of the entries for
		a null area or block, or a clear bit (a set bit in the
		complemented value). */

	}
	result := (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + int32(bits.TrailingZeros(uint(nword)))
	if result == math.MaxInt32 {
		return -1
	} else {
		return result
	}
}

/**
 *  Returns the index of the first bit that is set to <code>true</code> that
 *  occurs on or after the specified starting index. If no such it exists then
 *  -1 is returned.
 *  <p>
 *  To iterate over the <code>true</code> bits in a <code>SparseBitSet
 *  sbs</code>, use the following loop:
 *
 *  <pre>
 *  for( int i = sbbits.nextSetBit(0); i &gt;= 0; i = sbbits.nextSetBit(i+1) )
 *  {
 *      // operate on index i here
 *  }</pre>
 *
 * @param       i the index to start checking from (inclusive)
 * @return      the index of the next set bit
 * @exception   IndexOutOfBoundsException if the specified index is negative
 * @since       1.6
 */

func (this *BitSet) NextSetBit(i int32) int32 {
	/*  The index value (i) of this method is permitted to be Integer.MAX_VALUE,
	as this is needed to make the loop defined above work: just in case the
	bit labelled Integer.MAX_VALUE-1 is set. This case is not optimised:
	but eventually -1 will be returned, as this will be included with
	any search that goes off the end of the level1 array. */

	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	/*  This is the word from which the search begins. */
	w := i >> SHIFT3
	w3 := w & MASK3
	w2 := (w >> SHIFT2) & MASK2
	w1 := w >> SHIFT1

	word := wordType(0)
	aLength := int32(len(this.bits))

	/*  Is the next set bit in the same word at the nominated beginning bit
	(including the nominated beginning bit itself). The first check is
	whether the starting bit is within the structure at all. */
	var a2 b2DimType
	var a3 b1DimType
	if w1 < aLength {
		a2 = this.bits[w1]
		result := true
		if a2 != nil {
			if a3 = a2[w2]; a3 != nil {
				word = a3[w3] & wordType(^int64(0)<<uint(i))
				result = word == 0
			}
		}
		if result {
			/*  So now start a search though the rest of the entries for a bit. */
			w++
			w3 = w & MASK3
			w2 = (w >> SHIFT2) & MASK2
			w1 = w >> SHIFT1
		major:
			for ; w1 != aLength; w1++ {
				if a2 = this.bits[w1]; a2 != nil {
					for ; w2 != LENGTH2; w2++ {
						if a3 = a2[w2]; a3 != nil {
							for ; w3 != LENGTH3; w3++ {
								if word = a3[w3]; word != 0 {
									break major
								}
							}
							w3 = 0
						}
					}
					w2, w3 = 0, 0
				}
			}
		}
	}
	if w1 >= aLength {
		return -1
	} else {
		return (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + int32(bits.TrailingZeros(uint(word)))
	}
}

/**
 * Returns the index of the nearest bit that is set to {@code false}
 * that occurs on or before the specified starting index.
 * If no such bit exists, or if {@code -1} is given as the
 * starting index, then {@code -1} is returned.
 *
 * @param  i the index to start checking from (inclusive)
 * @return the index of the previous clear bit, or {@code -1} if there
 *         is no such bit
 * @throws IndexOutOfBoundsException if the specified index is less
 *         than {@code -1}
 * @since  1.2
 * @see java.util.BitSet#previousClearBit
 */
func (this *BitSet) PreviousClearBit(i int32) int32 {
	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}

	bits := this.bits
	aSize := int32(len(this.bits) - 1)

	w := i >> SHIFT3
	w3 := w & MASK3
	w2 := (w >> SHIFT2) & MASK2
	w1 := w >> SHIFT1
	if w1 > aSize {
		return i
	}

	if aSize < w1 {
		w1 = aSize
	}
	w4 := i % LENGTH4

	word := wordType(0)
	var a2 b2DimType
	var a3 b1DimType

	for ; w1 >= 0; w1-- {
		if a2 = bits[w1]; a2 == nil {
			return (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + w4
		}
		for ; w2 >= 0; w2-- {
			if a3 = a2[w2]; a3 == nil {
				return (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + w4
			}
			for ; w3 >= 0; w3-- {
				if word = a3[w3]; word == 0 {
					return (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + w4
				}
				for bitIdx := w4; bitIdx >= 0; bitIdx-- {
					if t := word & (1 << uint(bitIdx)); t == 0 {
						return (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + bitIdx
					}
				}
				w4 = LENGTH4_SIZE
			}
			w3 = LENGTH3_SIZE
		}
		w2 = LENGTH2_SIZE
	}
	return -1
}

/**
 * Returns the index of the nearest bit that is set to {@code true}
 * that occurs on or before the specified starting index.
 * If no such bit exists, or if {@code -1} is given as the
 * starting index, then {@code -1} is returned.
 *
 * @param  i the index to start checking from (inclusive)
 * @return the index of the previous set bit, or {@code -1} if there
 *         is no such bit
 * @throws IndexOutOfBoundsException if the specified index is less
 *         than {@code -1}
 * @since  1.2
 * @see java.util.BitSet#previousSetBit
 */

func (this *BitSet) PreviousSetBit(i int32) int32 {
	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	bits := this.bits
	aSize := int32(len(this.bits) - 1)

	/*  This is the word from which the search begins. */
	w := i >> SHIFT3
	w1 := w >> SHIFT1
	var w2, w3, w4 int32
	/*  But if its off the end of the array, start from the very end. */
	if w1 > aSize {
		w1 = aSize
		w2 = LENGTH2_SIZE
		w3 = LENGTH3_SIZE
		w4 = LENGTH4_SIZE
	} else {
		w2 = (w >> SHIFT2) & MASK2
		w3 = w & MASK3
		w4 = i % LENGTH4
	}
	word := wordType(0)
	var a2 b2DimType
	var a3 b1DimType

	for ; w1 >= 0; w1-- {
		if a2 = bits[w1]; a2 != nil {
			for ; w2 >= 0; w2-- {
				if a3 = a2[w2]; a3 != nil {
					for ; w3 >= 0; w3-- {
						if word = a3[w3]; word != 0 {
							for bitIdx := w4; bitIdx >= 0; bitIdx-- {
								if t := word & (1 << uint(bitIdx)); t != 0 {
									return (((w1 << SHIFT1) + (w2 << SHIFT2) + w3) << SHIFT3) + bitIdx
								}
							}
						}
						w4 = LENGTH4_SIZE
					}
				}
				w3 = LENGTH3_SIZE
				w4 = LENGTH4_SIZE
			}
		}
		w2 = LENGTH2_SIZE
		w3 = LENGTH3_SIZE
		w4 = LENGTH4_SIZE
	}
	return -1
}

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
func (this *BitSet) OrBit(i int32, value bool) {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}
	if value {
		this.Set(i)
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
func (this *BitSet) OrRangeBitSet(i, j int32, b *BitSet) {
	this.setScanner(i, j, b, orStrategy)
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
func (this *BitSet) OrBitSet(b *BitSet) {
	this.setScanner(0, b.bitsLength, b, orStrategy)
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

/**
 *  Sets the bit at the specified index.
 *
 * @param       i a bit index
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE
 * @since       1.6
 */
func (this *BitSet) Set(i int32) {
	if (i + 1) < 1 {
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
	a3[(w & MASK3)] |= wordType(bits.RotateLeft(uint(1), int(i)))
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

/**
 *  Sets the bit at the specified index to <code>false</code>.
 *
 * @param       i a bit index.
 * @exception   IndexOutOfBoundsException if the specified index is negative
 *              or equal to Integer.MAX_VALUE.
 * @since       1.6
 */
func (this *BitSet) Clear(i int32) {
	/*  In the interests of speed, no check is made here on whether the
	level3 block goes to all zero. This may be found and corrected
	in some later operation. */
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	if i >= this.bitsLength {
		return
	}
	w := i >> SHIFT3
	a2 := this.bits[w>>SHIFT1]
	if a2 == nil {
		return
	}
	a3 := a2[(w>>SHIFT2)&MASK2]
	if a3 == nil {
		return
	}
	a3[(w & MASK3)] &= ^wordType(int64(1) << uint(i)) //  Clear the indicated bit
	this.cache.hash = 0                               //  Invalidate size, etc.,
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
func (this *BitSet) ClearRange(i, j int32) {
	this.setScanner(i, j, nil, clearStrategy)
}

/**
 *  Sets all of the bits in this <code>SparseBitSet</code> to
 *  <code>false</code>.
 *
 * @since       1.6
 */
func (this *BitSet) clearAll() {
	/*  This simply resets to null all the entries in the set. */
	this.nullify(0)
}

/**
 *  Cloning this <code>SparseBitSet</code> produces a new
 *  <code>SparseBitSet</code> that is <i>equal</i>() to it. The clone of the
 *  bit set is another bit set that has exactly the same bits set to
 *  <code>true</code> as this bit set.
 *  <p>
 *  Note: the actual space allocated to the clone tries to minimise the actual
 *  amount of storage allocated to hold the bits, while still trying to
 *  keep access to the bits being a rapid as possible. Since the space
 *  allocated to a <code>SparseBitSet</code> is not normally decreased,
 *  replacing a bit set by its clone may be a way of both managing memory
 *  consumption and improving the rapidity of access.
 *
 * @return      a clone of this SparseBitSet
 * @since       1.6
 */
//public SparseBitSet clone()
func (this *BitSet) clone() (result *BitSet) {
	result = new(BitSet)

	reflect.Copy(reflect.ValueOf(*result), reflect.ValueOf(*this))

	/*  Clear out the shallow copy of the set array (which contains just
	copies of the references from this set), and then replace these
	by a deep copy (created by a "copy" from the set being cloned . */
	result.bits = nil
	result.resize(1)
	/*  Ensure the clone is not sharing a copy of a spare block with
	the cloned set, nor the cache set, nor any of the visitors (which
	are linked to their parent object) (Not all visitors actually use
	this link to their containing object, but they are reset here just
	in case of  future changes). */
	result.constructorHelper()
	result.equalsStrategy = nil
	result.setScanner(0, this.bitsLength, this, copyStrategy)
	return result
}

func highestOneBit(x int32) (result int32) {
	l := bits.Len(uint(x))
	if result == 0 {
		return
	}
	result = 1 << (uint(l) - 1)
	return
}

func (this *BitSet) resize(index int32) {
	/*  Find an array size that is a power of two that is as least as large
	enough to contain the index requested. */
	w1 := int32((index >> SHIFT3) >> SHIFT1)
	newSize := int32(highestOneBit(w1))
	if newSize == 0 {
		newSize = 1
	}
	if w1 >= newSize {
		newSize <<= 1
	}
	if newSize > MAX_LENGTH1 {
		newSize = MAX_LENGTH1
	}

	aLength1 := int32(0)
	if this.bits != nil {
		aLength1 = int32(len(this.bits))
	}

	if newSize != aLength1 || this.bits == nil {
		// only if the size needs to be changed
		temp := make(b3DimType, newSize) //  Get the new array
		if aLength1 != 0 {
			/*  If it exists, copy old array to the new array. */
			copy(temp, this.bits)
			this.nullify(0) //  Don't leave unused pointers around. */
		}
		this.bits = temp                //  Set new array as the set array
		this.bitsLength = math.MaxInt32 //  Index of last possible bit, plus one.
		if newSize != MAX_LENGTH1 {
			this.bitsLength = newSize * UNIT
		}
	}
}

/**
 *  Returns a string representation of this bit set. For every index for which
 *  this <code>SparseBitSet</code> contains a bit in the set state, the decimal
 *  representation of that index is included in the result. Such indices are
 *  listed in order from lowest to highest. If there is a subsequence of set
 *  bits longer than the value given by toStringCompaction, the subsequence
 *  is represented by the value for the first and the last values, with ".."
 *  between them. The individual bits, or the representation of sub-sequences
 *  are separated by ",&nbsp;" (a comma and a space) and surrounded by braces,
 *  resulting in a compact string showing (a variant of) the usual mathematical
 *  notation for a set of integers.
 *  <br>
 *  Example (with the default value of 2 for subsequences):
 *  <pre>
 *      SparseBitSet drPepper = new SparseBitSet();
 *  </pre>
 *  Now <code>drPepper.toString()</code> returns "<code>{}</code>".
 *  <br>
 *  <pre>
 *      drPepper.set(2);
 *  </pre>
 *  Now <code>drPepper.toString()</code> returns "<code>{2}</code>".
 *  <br>
 *  <pre>
 *      drPepper.set(3, 4);
 *      drPepper.set(10);
 *  </pre>
 *  Now <code>drPepper.toString()</code> returns "<code>{2..4, 10}</code>".
 *  <br>
 *  This method is intended for diagnostic use (as it is relatively expensive
 *  in time), but can be useful in interpreting problems in an application's use
 *  of a <code>SparseBitSet</code>.
 *
 * @return      a String representation of this SparseBitSet
 * @see         #toStringCompaction(int length)
 * @since       1.6
 */
func (this BitSet) String() string {
	var p string = "{"
	i := this.NextSetBit(0)
	/*  Loop so long as there is another bit to append to the String. */
	for i >= 0 {
		/*  Append that next bit */
		p += fmt.Sprintf("%v", i)
		/*  Find the position of the next bit to show. */
		j := this.NextSetBit(i + 1)
		if this.compactionCount > 0 {
			/*  Give up if there is no next bit to show. */
			if j < 0 {
				break
			}
			/*  Find the next clear bit is after the current bit, i.e., i */
			last := this.NextClearBit(i)
			/*  Compute the position of the next clear bit after the current
			subsequence of set bits. */
			if last < 0 {
				last = math.MaxInt32
			}
			/*  If the subsequence is more than the specified bits long, then
			collapse the subsequence into one entry in the String. */
			if (i + this.compactionCount) < last {
				p += fmt.Sprintf("..%v", last-1)
				/*  Having accounted for a subsequence of bits that are all set,
				recompute the label of the next bit to show. */
				j = this.NextSetBit(last)
			}
		}
		/*  If there is another set bit, put a comma and a space after the
		last entry in the String.  */
		if j >= 0 {
			p += ","
		}
		/*  Transfer to i the index of the next set bit. */
		i = j
	}
	/*  Terminate the representational String, and return it. */
	p += "}"
	return p
}

//==============================================================================
//      Internal methods
//==============================================================================

/**
 *  Throw the exception to indicate a range error. The <code>String</code>
 *  constructed reports all the possible errors in one message.
 *
 * @param       i lower bound for a operation
 * @param       j upper bound for a operation
 * @exception   IndexOutOfBoundsException indicating the range is not valid
 * @since       1.6
 */
func (this *BitSet) throwIndexOutOfBoundsException(i, j uint32) {
	var s string
	if i < 0 {
		s += fmt.Sprintf("(i=%v) < 0", i)
	} else if i == uint32(math.MaxInt32) {
		s += fmt.Sprintf("(i=%v)", i)
	}
	if j < 0 {
		if len(s) > 0 {
			s += ", "
		}
		s = s + fmt.Sprintf("(j=%v)<0", j)
	}
	if i > j {
		if len(s) > 0 {
			s += ", "
		}
		s += fmt.Sprintf("(i=%v) > (j=%v)", i, j)
	}
	panic(s)
}

func (this *BitSet) nullify(start int32) {
	aLength := int32(len(this.bits))
	if start < aLength {
		for w := start; w != aLength; w++ {
			this.bits[w] = nil
		}
		this.cache.hash = 0 //  Invalidate size, etc., values
	}
}

/**
 *  Intializes all the additional objects required for correct operation.
 *
 * @since       1.6
 */
func (this *BitSet) constructorHelper() {
	this.spare = make(b1DimType, LENGTH3)
	this.cache = new(cacheType)
	this.updateStrategy = new(updateStrategyType)
	this.updateStrategy.cache = this.cache
}

/**
 *  Scans over the bit set (and a second bit set if part of the operation) are
 *  all performed by this method. The properties and the operation executed
 *  are defined by a given strategy, which must be derived from the
 *  <code>AbstractStrategy</code>. The strategy defines how to operate on a
 *  single word, and on whole words that may or may not constitute a full
 *  block of words.
 *
 * @param       i the bit (inclusive) at which to start the scan
 * @param       j the bit (exclusive) at which to stop the scan
 * @param       b a SparseBitSet, if needed, the second SparseBitSet in the
 *              operation
 * @param       op the AbstractStrategy class defining the operation to be
 *              executed
 * @exception   IndexOutOfBoundsException
 * @since       1.6
 * @see         AbstractStrategy
 */
func (this *BitSet) setScanner(i, j int32, b *BitSet, op strateger) {

	/*  This method has been assessed as having a McCabe cyclomatic
	complexity of 47 (i.e., impossibly high). However, given that this
	method incorporates all the set scanning logic for all methods
	(with the exception of nextSetBit and nextClearBit, which themselves
	have high cyclomatic complexities of 13), and is attempting to minimise
	execution time (hence deals with processing shortcuts), it cannot be
	expected to be simple. In fact, the work of lining up level3 blocks
	proceeds step-wise, and each sub-section piece is reasonably
	straight-forward. Nevertheless, the number of paths is high, and
	caution is advised in attempting to correct anything. */

	/*  Do whatever the strategy needs to get started, and do whatever initial
	checking is needed--fail here if needed before much else is done. */
	if op.start(b) {
		this.cache.hash = 0
	}

	if j < i || (i+1) < 1 {
		panic(fmt.Sprintf("throwIndexOutOfBoundsException(%v,%v)", i, j))
	}

	if i == j {
		return
	}

	/*  Get the values of all the short-cut options. */
	properties := op.properties()
	f_op_f_eq_f := (properties & F_OP_F_EQ_F) != 0
	f_op_x_eq_f := (properties & F_OP_X_EQ_F) != 0
	x_op_f_eq_f := (properties & X_OP_F_EQ_F) != 0
	x_op_f_eq_x := (properties & X_OP_F_EQ_X) != 0

	/*  Index of the current word, and mask for the first word,
	to be processed in the bit set. */
	u := int32(i) >> SHIFT3
	//final long um = ~0L << i;
	um := wordType(int64(^0) << uint(i))

	/*  Index of the final word, and mask for the final word,
	to be processed in the bit set. */
	v := int32((j - 1)) >> SHIFT3
	// final long vm = ~0L >>> -j;
	vm := ^wordType(int64(^0) << uint(j))

	/*  Set up the two bit arrays (if the second exists), and their
	corresponding lengths (if any). */
	a1 := this.bits //  Level1, i.e., the bit arrays
	aLength1 := int32(len(this.bits))

	var b1 b3DimType
	var bLength1 int32 = 0

	if b1 != nil {
		b1 = b.bits
		bLength1 = int32(len(b.bits))
	}

	/*  Calculate the initial values of the parts of the words addresses,
	as well as the location of the final block to be processed.  */
	u1 := u >> SHIFT1
	u2 := (u >> SHIFT2) & MASK2
	u3 := u & MASK3
	v1 := v >> SHIFT1
	v2 := (v >> SHIFT2) & MASK2
	v3 := v & MASK3
	lastA3Block := (v1 << LEVEL2) + v2

	/*  Initialize the local copies of the counts of blocks and areas; and
	whether there is a partial first block.  */
	var a2CountLocal int32 = 0
	var a3CountLocal int32 = 0
	notFirstBlock := u == 0 && um == ^wordType(0)

	/*  The first level2 is cannot be judged empty if not being scanned from
	the beginning. */
	a2IsEmpty := u2 == 0 //  Presumption
	for i < j {
		/*  Determine if there is a level2 area in both the a and the b set,
		and if so, set the references to these areas. */
		var a2 b2DimType

		haveA2 := u1 < aLength1
		if haveA2 {
			a2 = a1[u1]
			haveA2 = haveA2 && a2 != nil
		}

		var b2 b2DimType
		haveB2 := u1 < bLength1 && b1 != nil
		if haveB2 {
			b2 = b1[u1]
			haveB2 = haveB2 && b2 != nil
		}
		/*  Handling of level 2 empty areas: determined by the
		properties of the strategy. It is necessary to actually visit
		the first and last blocks of a scan, since not all of the block
		might participate in the operation, hence making decision based
		on just the references to the blocks could be wrong. */
		if notFirstBlock &&
			u1 != v1 &&
			(!haveA2 && !haveB2 && f_op_f_eq_f ||
				!haveA2 && f_op_x_eq_f ||
				!haveB2 && x_op_f_eq_f) {
			//nested if!
			if u1 < aLength1 {
				a1[u1] = nil
			}
		} else {
			limit2 := LENGTH2
			if u1 == v1 {
				limit2 = int32(v2 + 1)
			}

			for u2 != int32(limit2) {
				/*  Similar logic applied here as for the level2 blocks.
				The initial and final block must be examined. In other
				cases, it may be possible to make a decision based on
				the value of the references, as indicated by the
				properties of the strategy. */
				a3IsSpare := false

				var a3 b1DimType

				haveA3 := haveA2
				if haveA3 {
					a3 = a2[u2]
					haveA3 = haveA3 && a3 != nil
				}

				var b3 b1DimType
				haveB3 := haveB2
				if haveB3 {
					b3 = b2[u2]
					haveB3 = haveB3 && b3 != nil
				}

				a3Block := (u1 << LEVEL2) + u2
				notLastBlock := lastA3Block != a3Block

				/*  Handling of level 3 empty areas: determined by the
				properties of the strategy. */
				if (!haveA3 && !haveB3 && f_op_f_eq_f || !haveA3 && f_op_x_eq_f || !haveB3 && x_op_f_eq_f) && notFirstBlock && notLastBlock {
					/*  Do not need level3 block, so remove it, and move on. */
					if haveA2 {
						a2[u2] = nil
					}
				} else {
					/*  So what is needed is the level3 block. */
					base3 := a3Block << SHIFT2
					limit3 := LENGTH3
					if !notLastBlock {
						limit3 = int32(v3)
					}
					if !haveA3 {
						a3 = this.spare
						a3IsSpare = true
					}
					if !haveB3 {
						b3 = ZERO_BLOCK
					}
					isZero := false
					if notFirstBlock && notLastBlock {
						if x_op_f_eq_x && !haveB3 {
							isZero = isZeroBlock(a3)
						} else {
							isZero = op.block(base3, 0, LENGTH3, a3, b3)
						}
					} else {
						/*  Partial block to process. */
						if notFirstBlock {
							/*  By implication, this is the last block */
							isZero = op.block(base3, 0, limit3, a3, b3)
							//  Do the whole words
							isZero = isZero && op.word(base3, limit3, a3, b3, vm)
							//  And then the final word
						} else {
							// u, v are correct if first block
							if u == v { //  Scan starts and ends in one word
								isZero = op.word(base3, u3, a3, b3, (um & vm))
							} else {
								// Scan starts in this a3 block
								isZero = op.word(base3, u3, a3, b3, um)
								//  First word
								isZero = isZero && op.block(base3, u3+1, limit3, a3, b3)
								//  Remainder of full words in block
								if limit3 != LENGTH3 {
									isZero = isZero && op.word(base3, limit3, a3, b3, vm)
								}
								//  If there is a partial word left
							}
							notFirstBlock = true //  Only one first block
						}
						if isZero {
							isZero = isZeroBlock(a3)
						}
						// If not known to have a non-zero
						// value, be sure whether all zero.
					}
					if isZero { //  The resulting a3 block has no values
						// nested if!
						/*  If there is an level 2 area make the entry for this
						level3 block be a null (i.e., remove any a3 block ). */
						if haveA2 {
							a2[u2] = nil
						}
					} else {
						/*  If the a3 block used was the spare block, put it
						into current level2 area; get a new spare block. */
						if a3IsSpare {
							if int32(i) >= this.bitsLength { //Check that the set is large
								//  enough to take the new block
								this.resize(i) //  Make it large enough
								a1 = this.bits //  Update reference and length
								aLength1 = int32(len(a1))
							}
							if a2 == nil { //  Ensure a level 2 area
								a2 = make(b2DimType, LENGTH2)
								a1[u1] = a2
								haveA2 = true //  Ensure know level2 not empty
							}
							a2[u2] = a3 //  Insert the level3 block
							a3IsSpare = false
							this.spare = make(b1DimType, LENGTH3) // Replace the spare

						}
						a3CountLocal++ // Count the level 3 block
					}
					a2IsEmpty = a2IsEmpty && !(haveA2 && a2[u2] != nil)
				} //  Keep track of level 2 usage
				u2++
				u3 = 0
			} /* end while ( u2 != limit2 ) */
			/*  If the loop finishes without completing the level 2, it may
			be left with a reference but still be all null--this is OK. */
			if u2 == LENGTH2 && a2IsEmpty && u1 < aLength1 {
				a1[u1] = nil
			} else {
				a2CountLocal++ //  Count level 2 areas
			}
		}
		/*  Advance the value of u based on what happened. */
		u1++
		u = (u1 << SHIFT1)
		i = u << SHIFT3
		u2 = 0 //  u3 = 0
		//  Compute next word and bit index
		if i < 0 {
			i = math.MaxInt32 //  Don't go over the end
		}

	} /* end while( i < j ) */

	/*  Do whatever the strategy needs in order to finish. */
	op.finish(a2CountLocal, a3CountLocal)
}

/**
 *  The entirety of the bit set is examined, and the various statistics of
 *  the bit set (size, length, cardinality, hashCode, etc.) are computed. Level
 *  arrays that are empty (i.e., all zero at level 3, all null at level 2) are
 *  replaced by null references, ensuring a normalized representation.
 *
 * @since       1.6
 */
func (this *BitSet) statisticsUpdate() {
	if this.cache.hash != 0 {
		return
	}
	this.setScanner(0, this.bitsLength, nil, this.updateStrategy)
}

/**
 *  Returns the number of bits set to <code>true</code> in this
 *  <code>SparseBitSet</code>.
 *
 * @return      the number of bits set to true in this SparseBitSet
 * @since       1.6
 */

func (this *BitSet) Cardinality() int32 {
	this.statisticsUpdate() // Update size, cardinality and length values
	return this.cache.cardinality
}

/**
 *  Convenience method for statistics if the individual results are not needed.
 *
 * @return      a String detailing the statistics of the bit set
 * @see         #statistics(String[])
 * @since       1.6
 */
func (this *BitSet) StatisticsAll() string {
	return this.Statistics(nil)
}

/**
 *  Determine, and create a String with the bit set statistics. The statistics
 *  include: Size, Length, Cardinality, Total words (<i>i.e.</i>, the total
 *  number of 64-bit "words"), Set array length (<i>i.e.</i>, the number of
 *  references that can be held by the top level array, Level2 areas in use,
 *  Level3 blocks in use,, Level2 pool size, Level3 pool size, and the
 *  Compaction count.
 *  <p>
 *  This method is intended for diagnostic use (as it is relatively expensive
 *  in time), but can be useful in understanding an application's use of a
 *  <code>SparseBitSet</code>.
 *
 * @param       values an array for the individual results (if not null)
 * @return      a String detailing the statistics of the bit set
 * @since       1.6
 */
func (this *BitSet) Statistics(values []string) string {
	this.statisticsUpdate() //  Ensure statistics are up-to-date
	v := make([]string, Statistics_Values_Length)
	/*  Assign the statistics values to the appropriate entry. The order
	of the assignments does not matter--the ordinal serves to get the
	values into the matching order with the labels from the enumeration. */
	v[Size] = strconv.Itoa(int(this.Size()))
	v[Length] = strconv.Itoa(int(this.Length()))
	v[Cardinality] = strconv.Itoa(int(this.Cardinality()))
	v[Total_words] = strconv.Itoa(int(this.cache.count))
	v[Set_array_length] = strconv.Itoa(len(this.bits))
	v[Set_array_max_length] = strconv.Itoa(int(MAX_LENGTH1))
	v[Level2_areas] = strconv.Itoa(int(this.cache.a2Count))
	v[Level2_area_length] = strconv.Itoa(int(LENGTH2))
	v[Level3_blocks] = strconv.Itoa(int(this.cache.a3Count))
	v[Level3_block_length] = strconv.Itoa(int(LENGTH3))
	v[Compaction_count_value] = strconv.Itoa(int(this.compactionCount))

	/*  Determine the longest label, so that the equal signs may be lined-up. */
	for i := range values {
		if i < len(v) {
			values[i] = v[i]
		}
	}

	/*  Build a String that has for each statistic, the name of the statistic,
	padding, and equals sign, and the value. The "Load_factor_value",
	"Average_length_value", and "Average_chain_length" are printed as
	floating point values. */
	var kvs string
	for i, s := range values {
		st := StatisticsType(i)
		kvs = kvs + st.String() + " = " + s + "\n"
	}
	return kvs
}

//=============================================================================
//  Statistics enumeration
//=============================================================================

/**
 *  These enumeration values are used as labels for the values in the String
 *  created by the <i>statistics</i>() method. The values of the corresponding
 *  statistics are <code>int</code>s, except for the loadFactor and
 *  Average_chain_length values, which are <code>float</code>s.
 *  <p>
 *  An array of <code>String</code>s may be obtained containing a
 *  representation of each of these values. An element of such an array, say,
 *  <code>values</code>, may be accessed, for example, by:
 *  <pre>
 *      values[SparseBitSet.statistics.Buckets_available.ordinal()]</pre>
 *
 * @see         #statistics(String[])
 */
type StatisticsType int

const (
	/**
	 *  The size of the bit set, as give by the <i>size</i>() method.
	 */
	Size StatisticsType = iota // 0
	/**
	 *  The length of the bit set, as give by the <i>length</i>() method.
	 */
	Length // 1
	/**
	 *  The cardinality of the bit set, as give by the <i>cardinality</i>() method.
	 */
	Cardinality // 2
	/**
	 *  The total number of non-zero 64-bits "words" being used to hold the
	 *  representation of the bit set.
	 */
	Total_words // 3
	/**
	 *  The length of the bit set array.
	 */
	Set_array_length // 4
	/**
	 *  The maximum permitted length of the bit set array.
	 */
	Set_array_max_length // 5
	/**
	 *  The number of level2 areas.
	 */
	Level2_areas // 6
	/**
	 *  The length of the level2 areas.
	 */
	Level2_area_length // 7
	/**
	 *  The total number of level3 blocks in use.
	 */
	Level3_blocks // 8
	/**
	 *  The length of the level3 blocks.
	 */
	Level3_block_length // 9
	/**
	 *  Is the value that determines how the <i>toString</i>() conversion is
	 *  performed.
	 * @see         #toStringCompaction(int)
	 */
	Compaction_count_value // 10
	//
	Statistics_Values_Length
)

func (st StatisticsType) String() string {
	switch st {
	case Size:
		return "Size"
	case Length:
		return "Length"
	case Cardinality:
		return "Cardinality"
	case Total_words:
		return "Total-words"
	case Set_array_length:
		return "Set-array-length"
	case Set_array_max_length:
		return "Set-array-max-length"
	case Level2_areas:
		return "Level2-areas"
	case Level2_area_length:
		return "Level2-area-length"
	case Level3_blocks:
		return "Level3-blocks"
	case Level3_block_length:
		return "Level3-block-length"
	case Compaction_count_value:
		return "Compaction-count-value"
	default:
		panic(fmt.Sprintf("Unknown statistics value %v", st))
	}
}

//=============================================================================
//  A set of cached statistics values, recomputed when necessary
//=============================================================================

/**
 *  This class holds the values related to various statistics kept about the
 *  bit set. These values are not kept continuously up-to-date. Whenever the
 *  values become invalid, the field <i>hash</i> is set to zero, indicating
 *  that an update is required.
 *
 * @see     #statisticsUpdate()
 */
type cacheType struct {
	/**
	 *  <i>hash</i> is updated by the <i>statisticsUpdate</i>() method.
	 *  If the <i>hash</i> value is zero, it is assumed that <b><i>all</i></b>
	 *  the cached values are stale, and must be updated.
	 */
	hash uint64

	/**
	 *  <i>size</i> is updated by the <i>statisticsUpdate</i>() method.
	 *  If the <i>hash</i> value is zero, it is assumed the all the cached
	 *  values are stale, and must be updated.
	 */
	size int32

	/**
	 *  <i>cardinality</i> is updated by the <i>statisticsUpdate</i>() method.
	 *  If the <i>hash</i> value is zero, it is assumed the all the cached
	 *  values are stale, and must be updated.
	 */
	cardinality int32

	/**
	 *  <i>length</i> is updated by the <i>statisticsUpdate</i>() method.
	 *  If the <i>hash</i> value is zero, it is assumed the all the cached
	 *  values are stale, and must be updated.
	 */
	length int32

	/**
	 *  <i>count</i> is updated by the <i>statisticsUpdate</i>() method.
	 *  If the <i>hash</i> value is zero, it is assumed the all the cached
	 *  values are stale, and must be updated.
	 */
	count int32

	/**
	 *  <i>a2Count</i> is updated by the <i>statisticsUpdate</i>()
	 *  method, and will only be correct immediately after a full update. The
	 *  <i>hash</i> value is must be zero for all values to be updated.
	 */
	a2Count int32

	/**
	 *  <i>a3Count</i> is updated by the <i>statisticsUpdate</i>() method,
	 *  and will only be correct immediately after a full update. The
	 *  <i>hash</i> value is must be zero for all values to be updated.
	 */
	a3Count int32
}

//=============================================================================
//  Strategies based on the Strategy super-class describing logical operations
//=============================================================================

type strateger interface {
	properties() int32
	start(*BitSet) bool
	word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool
	block(base, u3, v3 int32, a3, b3 b1DimType) bool
	finish(a2Count, a3Count int32)
}

/** If the operation requires that when matching level2 areas or level3
 *  blocks are null, that no action is required, then this property is
 *  required. Corresponds to the top-left entry in the logic diagram for the
 *  operation being 0. For all the defined actual logic operations ('and',
 *  'andNot', 'or', and 'xor', this will be true, because for all these,
 *  "false" op "false" = "false".
 */
const F_OP_F_EQ_F = 0x1

/** If when level2 areas or level3 areas from the this set are null will
 *  require that area or block to remain null, irrespective of the value of
 *  the matching structure from the other set, then this property is required.
 *  Corresponds to the first row in the logic diagram being all zeros. For
 *  example, this is true for 'and' as well as 'andNot', and for 'clear', since
 *  false" & "x" = "false", and "false" &! "x" = "false".
 */
const F_OP_X_EQ_F = 0x2

/** If when level2 areas or level3 areas from the other set are null will
 *  require the matching area or block in this set to be set to null,
 *  irrespective of the current values in the matching structure from the
 *  this, then this property is required. Corresponds to the first column
 *  in the logic diagram being all zero. For example, this is true for
 *  'and', since "x" & "false" = "false", as well as for 'clear'.
 */
const X_OP_F_EQ_F = 0x4

/** If when a level3 area from the other set is null will require the
 *  matching area or block in this set to be left as it is, then this property
 *  is required. Corresponds to the first column of the logic diagram being
 *  equal to the left hand operand column. For example, this is true for 'or',
 *  'xor', and 'andNot', since for all of these "x" op "false" = "x".
 */
const X_OP_F_EQ_X = 0x8

func isZeroBlock(a3 b1DimType) bool {
	for _, word := range a3 {
		if word != 0 {
			return false
		}
	}
	return true
}

/**
 *  And of two sets. Where the <i>a</i> set is zero, it remains zero (i.e.,
 *  without entries or with zero words). Similarly, where the <i>b</i> set is
 *  zero, the <i>a</i> becomes zero (i.e., without entries).
 *  <p>
 *  If level1 of the <i>a</i> set is longer than level1 of the bit set
 *  <i>b</i>, then the unmatched virtual "entries" of the <i>b</i> set (beyond
 *  the actual length of <i>b</i>) corresponding to these are all false, hence
 *  the result of the "and" operation will be to make all these entries in this
 *  set to become false--hence just remove them, and then scan only those
 *  entries that could match entries in the bit set<i>b</i>. This clearing of
 *  the remainder of the <i>a</i> set is accomplished by selecting both
 *  <i>F_OP_X_EQ_F</i> and <i>X_OP_F_EQ_F</i>.
 *
 *  <pre>
 *  and| 0 1
 *    0| 0 0
 *    1| 0 1 <pre>
 */
type andStrategyType struct{}

func (st andStrategyType) properties() int32 {
	return F_OP_F_EQ_F + F_OP_X_EQ_F + X_OP_F_EQ_F
}

func (st andStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	return true
}

func (st andStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3]&b3[u3] | ^mask
	return a3[u3] == 0
}

func (st andStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] & b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	return
}

func (st andStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  AndNot of two sets. Where the <i>a</i> set is zero, it remains zero
 *  (i.e., without entries or with zero words). On the other hand, where the
 *  <i>b</i> set is zero, the <i>a</i> remains unchanged.
 *  <p>
 *  If level1 of the <i>a</i> set is longer than level1 of the bit set
 *  <i>b</i>, then the unmatched virtual "entries" of the <i>b</i> set (beyond
 *  the actual length of <i>b</i>) corresponding to these are all false, hence
 *  the result of the "and" operation will be to make all these entries in this
 *  set to become false--hence just remove them, and then scan only those
 *  entries that could match entries in the bit set<i>b</i>. This clearing of
 *  the remainder of the <i>a</i> set is accomplished by selecting both
 *  <i>F_OP_X_EQ_F</i> and <i>X_OP_F_EQ_F</i>.
 *
 *  <pre>
 * andNot| 0 1
 *      0| 0 0
 *      1| 1 0 <pre>
 */
type andNotStrategyType struct{}

func (st andNotStrategyType) properties() int32 {
	return F_OP_F_EQ_F + F_OP_X_EQ_F + X_OP_F_EQ_X
}

func (st andNotStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	return true
}

func (st andNotStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] & ^(b3[u3] & mask)
	return a3[u3] == 0
}

func (st andNotStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] & ^b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st andNotStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Clear clears bits in the <i>a</i> set.
 *
 * <pre>
 * clear| 0 1
 *     0| 0 0
 *     1| 0 0 <pre>
 */
type clearStrategyType struct{}

func (st clearStrategyType) properties() int32 {
	return F_OP_F_EQ_F + F_OP_X_EQ_F
}

func (st clearStrategyType) start(b *BitSet) bool {
	return true
}

func (st clearStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] & ^mask
	return a3[u3] == 0
}

func (st clearStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	if u3 != 0 || v3 != LENGTH3 {
		for w3 := u3; w3 != v3; w3 = w3 + 1 {
			a3[w3] = 0
		}
	}
	return true
}
func (st clearStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Copies the needed parts of the <i>b</i> set to the <i>a</i> set.
 *
 * <pre>
 * get| 0 1
 *   0| 0 1
 *   1| 0 1 <pre>
 */
type copyStrategyType struct{}

func (st copyStrategyType) properties() int32 {
	return F_OP_F_EQ_F + X_OP_F_EQ_F
}

func (st copyStrategyType) start(b *BitSet) bool {
	return true
}

func (st copyStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = b3[u3] & mask
	return a3[u3] == 0
}

func (st copyStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st copyStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Equals compares bits in the <i>a</i> set with those in the <i>b</i> set.
 *  None of the values in either set are changed, although the <i>a</i> set
 *  may have all zero level 3 blocks replaced by null references (and
 *  similarly at level 2).
 *
 * <pre>
 * equals| 0 1
 *      0| 0 -
 *      1| - - <pre>
 */
type equalsStrategyType struct {
	result bool
}

func (st equalsStrategyType) properties() int32 {
	return F_OP_F_EQ_F
}

func (st *equalsStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	st.result = true
	return false
	/*  Equals does not change the content of the set, hence hash need
	not be reset. */
}

func (st *equalsStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	word := a3[u3]
	st.result = st.result && ((word & mask) == (b3[u3] & mask))
	return word == 0
}

func (st *equalsStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		word := a3[w3]
		st.result = st.result && word == b3[w3]
		isZero = isZero && word == 0
	}
	return
}

func (st equalsStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Flip inverts the bits of the <i>a</i> set within the given range.
 *
 * <pre>
 * flip| 0 1
 *    0| 1 1
 *    1| 0 0 <pre>
 */
type flipStrategyType struct{}

func (st flipStrategyType) properties() int32 {
	return 0
}

func (st flipStrategyType) start(b *BitSet) bool {
	return true
}

func (st flipStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] ^ mask
	return a3[u3] == 0
}

func (st flipStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] ^ wordType(^wordType(0))
		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st flipStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Intersect has a true result if any word in the <i>a</i> set has a bit
 *  in common with the <i>b</i> set. During the scan of the <i>a</i> set
 *  blocks (and areas) that are all zero may be replaced with empty blocks
 *  and areas (null references), but the value of the set is not changed
 *  (which is why X_OP_F_EQ_F is not selected, since this would cause
 *  parts of the <i>a</i> set to be zero-ed out).
 *
 * <pre>
 * intersect| 0 1
 *         0| 0 0
 *         1| 1 1 <pre>
 */
type intersectsStrategyType struct {
	result bool
}

func (st intersectsStrategyType) properties() int32 {
	return F_OP_F_EQ_F + F_OP_X_EQ_F
}

func (st *intersectsStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	st.result = false
	return false
	/*  Intersect does not change the content of the set, hence hash need
	not be reset. */
}

func (st *intersectsStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	word := a3[u3]
	st.result = st.result || ((word & b3[u3] & mask) != 0)
	return word == 0
}

func (st *intersectsStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		word := a3[w3]
		st.result = st.result || ((word & b3[w3]) != 0)
		isZero = isZero && word == 0
	}
	return
}
func (st intersectsStrategyType) finish(a2Count, a3Count int32) {}

/**
 *  Or of two sets. Where the <i>a</i> set is one, it remains one. Similarly,
 *  where the <i>b</i> set is one, the <i>a</i> becomes one. If both sets have
 *  zeros in corresponding places, a zero results. Whole blocks or areas that
 *  are or become zero are replaced by null arrays.
 *  <p>
 *  If level1 of the <i>a</i> set is longer than level1 of the bit set
 *  <i>b</i>, then the unmatched entries of the <i>a</i> set (beyond
 *  the actual length of <i>b</i>) corresponding to these remain unchanged. *
 *  <pre>
 *   or| 0 1
 *    0| 0 1
 *    1| 1 1 <pre>
 */

type orStrategyType struct{}

func (st orStrategyType) properties() int32 {
	return F_OP_F_EQ_F + X_OP_F_EQ_X
}

func (st orStrategyType) start(b *BitSet) bool {
	return true
}

func (st orStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] | (b3[u3] & mask)
	return a3[u3] == 0
}

func (st orStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] | b3[w3]

		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st orStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Set creates entries everywhere within the range. Hence no empty level2
 *  areas or level3 blocks are ignored, and no empty (all zero) blocks are
 *  returned.
 *
 *  <pre>
 * set| 0 1
 *   0| 1 1
 *   1| 1 1 <pre>
 */

type setStrategyType struct{}

func (st setStrategyType) properties() int32 {
	return 0
}

func (st setStrategyType) start(b *BitSet) bool {
	return true
}

func (st setStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] | mask
	return a3[u3] == 0
}

func (st setStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = ^wordType(0)
	}
	isZero = false
	return
}
func (st setStrategyType) finish(a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Update the seven statistics that are computed for each set. These are
 *  updated by calling <i>statisticsUpdate</i>, which uses this strategy.
 *
 *  <pre>
 *  update| 0 1
 *       0| 0 0
 *       1| 1 1 <pre>
 *
 * @see SparseBitSet#statisticsUpdate()
 */
type updateStrategyType struct {
	/**
	 *  Working space for find the size and length of the bit set. Holds the
	 *  index of the first non-empty word in the set.
	 */
	wMin int32

	/**
	 *  Working space for find the size and length of the bit set. Holds copy of
	 *  the first non-empty word in the set.
	 */
	wordMin wordType

	/**
	 *  Working space for find the size and length of the bit set. Holds the
	 *  index of the last non-empty word in the set.
	 */
	wMax int32

	/**
	 *  Working space for find the size and length of the bit set. Holds a copy
	 *  of the last non-empty word in the set.
	 */
	wordMax wordType

	/**
	 *  Working space for find the hash value of the bit set. Holds the
	 *  current state of the computation of the hash value. This value is
	 *  ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	hash uint64

	/**
	 *  Working space for keeping count of the number of non-zero words in the
	 *  bit set. Holds the current state of the computation of the count. This
	 *  value is ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	count int32

	/**
	 *  Working space for counting the number of non-zero bits in the bit set.
	 *  Holds the current state of the computation of the cardinality.This
	 *  value is ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	cardinality int32

	cache *cacheType
}

func (st updateStrategyType) properties() int32 {
	return F_OP_F_EQ_F + F_OP_X_EQ_F
}

func (st *updateStrategyType) start(b *BitSet) bool {
	st.hash = 1234     // Magic number
	st.wMin = -1       // index of first non-zero word
	st.wordMin = 0     // word at that index
	st.wMax = 0        // index of last non-zero word
	st.wordMax = 0     // word at that index
	st.count = 0       // count of non-zero words in whole set
	st.cardinality = 0 // count of non-zero bits in the whole set
	return false
}

func (st updateStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	word := a3[u3]
	word1 := word & mask
	if word1 != 0 {
		st.compute(base+u3, word1)
	}
	return word == 0
}

func (st *updateStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true //  Presumption
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		if word := a3[w3]; word != 0 {
			isZero = false
			st.compute(base+w3, word)
		}
	}
	isZero = false
	return
}

func (st *updateStrategyType) finish(a2Count, a3Count int32) {
	st.cache.a2Count = a2Count
	st.cache.a3Count = a3Count
	st.cache.count = st.count
	st.cache.cardinality = st.cardinality
	st.cache.length = (st.wMax+1)*LENGTH4 - int32(bits.LeadingZeros(uint(st.wordMax)))
	st.cache.size = st.cache.length - st.wMin*LENGTH4 - int32(bits.LeadingZeros(uint(st.wordMin)))
	st.cache.hash = ((st.hash >> INTEGER_SIZE) ^ st.hash)
}

func (st *updateStrategyType) compute(index int32, word wordType) {
	/*  Count the number of actual words being used. */
	st.count++
	/*  Continue to accumulate the hash value of the set. */
	st.hash = st.hash ^ (uint64(word) * uint64(index+1))
	/*  The first non-zero word contains the first actual bit of the
	    set. The location of this bit is used to compute the set size. */
	if st.wMin < 0 {
		st.wMin = index
		st.wordMin = word
	}
	/*  The last non-zero word contains the last actual bit of the set.
	    The location of this bit is used to compute the set length. */
	st.wMax = index
	st.wordMax = word
	/*  Count the actual bits, so as to get the cardinality of the set. */
	st.cardinality = st.cardinality + int32(bits.OnesCount(uint(word)))
}

//-----------------------------------------------------------------------------
/**
 *  The XOR of level3 blocks is computed.
 *
 * <pre>
 * xor| 0 1
 *   0| 0 1
 *   1| 1 0 <pre>
 */
type xorStrategyType struct{}

func (st xorStrategyType) properties() int32 {
	return F_OP_F_EQ_F + X_OP_F_EQ_X
}
func (st xorStrategyType) start(b *BitSet) bool {
	return true
}

func (st xorStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] ^ (b3[u3] & mask)
	return a3[u3] == 0
}

func (st xorStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] ^ b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	isZero = false
	return
}
func (st xorStrategyType) finish(a2Count, a3Count int32) {}
