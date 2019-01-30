package sparse

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
)

type wordType uint64
type b3DimType [][][]wordType
type b2DimType [][]wordType
type b1DimType []wordType

/** An empty level 3 block is kept for use when scanning. When a source block
 *  is needed, and there is not already one in the corresponding bit set, the
 *  ZERO_BLOCK is used (as a read-only block). It is a source of zero values
 *  so that code does not have to test for a null level3 block. This is a
 *  static block shared everywhere.
 */
var iZeroBlock = make(b1DimType, cLength3)

//BitSet ....
type BitSet struct {
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
	 *  A spare level 3 block is kept for use when scanning. When a target block
	 *  is needed, and there is not already one in the bit set, the spare is
	 *  provided. If non-zero values are placed into this block, it is moved to the
	 *  resulting set, and a new spare is acquired. Note: a new spare needs to
	 *  be allocated when the set is cloned (so that the spare is not shared
	 *  between two sets).
	 */
	spare b1DimType
	/**
	 *  This value controls for format of the toString() output.
	 * @see #toStringCompaction(int)
	 */
	compactionCount int32

	/**
	 *  For the current size of the bits array, this is the maximum possible
	 *  length of the bit set, i.e., the index of the last possible bit, plus one.
	 *  Note: this not the value returned by <i>length</i>().
	 * @see #resize(int)
	 * @see #length()
	 */
	bitsLength int32

	/**
	 *  Word and block <b>equals</b> strategy.
	 */
	//equalsStrategy *equalsStrategyType
	/**
	 *  Word and block <b>update</b> strategy.
	 */
	//updateStrategy *updateStrategyType

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
	/**
	 *  Holds reference to the cache of statistics values computed by the
	 *  UpdateStrategy
	 * @see SparseBitSet.Cache
	 * @see SparseBitSet.UpdateStrategy
	 */
	cache cacheType
}

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
func (bs *BitSet) GetBit(i int32) bool {
	if (i + 1) < 1 {
		panic(fmt.Sprintf("IndexOutOfBoundsException: i=%v", i))
	}

	w := i >> cShift3
	if i < bs.bitsLength {
		a2 := bs.bits[w>>cShift1]
		if a2 != nil {
			a3 := a2[(w>>cShift2)&(cMask2)]
			if a3 != nil {
				result := a3[(w&cMask3)] & (wordType(uint(1) << remainderOf64(i)))
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

//GetBitSetFromRange ...
func (bs *BitSet) GetBitSetFromRange(i, j int32) *BitSet {
	result := newWithSizeAndCompactionCount(j, bs.compactionCount)
	result.setScanner(i, j, bs, copyStrategy)
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
func (bs *BitSet) IntersectsRangeBitSet(i, j int32, b *BitSet) bool {
	s := new(intersectsStrategyType)
	bs.setScanner(i, j, b, s)
	return s.result
}

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
// public boolean intersects(SparseBitSet b)
func (bs *BitSet) IntersectsBitSet(b *BitSet) bool {
	bmax := bs.bitsLength
	if b.bitsLength > bmax {
		bmax = b.bitsLength
	}
	s := new(intersectsStrategyType)
	bs.setScanner(0, bmax, b, s)
	return s.result
}

/**
 *  Returns true if this <code>SparseBitSet</code> contains no bits that are
 *  set to <code>true</code>.
 *
 * @return      the boolean indicating whether this SparseBitSet is empty
 * @since       1.6
 */
func (bs *BitSet) IsEmpty() bool {
	bs.statisticsUpdate()
	return bs.cache.cardinality == 0
}

/**
 *  Returns the "logical length" of this <code>SparseBitSet</code>: the index
 *  of the highest set bit in the <code>SparseBitSet</code> plus one. Returns
 *  zero if the <code>SparseBitSet</code> contains no set bits.
 *
 * @return      the logical length of this SparseBitSet
 * @since       1.6
 */
func (bs *BitSet) Length() int32 {
	bs.statisticsUpdate()
	return bs.cache.length
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
func (bs *BitSet) Size() int32 {
	bs.statisticsUpdate()
	return bs.cache.size
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
func (bs *BitSet) NextClearBit(i int32) int32 {
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
	w := i >> cShift3
	w3 := w & cMask3
	w2 := (w >> cShift2) & cMask2
	w1 := w >> cShift1

	nword := wordType(^uint64(0) << remainderOf64(i))
	aLength := int32(len(bs.bits))

	/*  Is the next clear bit in the same word at the nominated beginning bit
	(including the nominated beginning bit itself). The first check is
	whether the starting bit is within the structure at all. */

	if w1 < aLength {
		if a2 := bs.bits[w1]; a2 != nil {
			if a3 := a2[w2]; a3 != nil {
				if nword = ^a3[w3] & wordType(^uint64(0)<<remainderOf64(i)); nword == 0 {
					w++
					w3 = w & cMask3
					w2 = (w >> cShift2) & cMask2
					w1 = w >> cShift1
					nword = ^wordType(0)
				loop:
					for ; w1 != aLength; w1++ {
						if a2 = bs.bits[w1]; a2 == nil {
							break
						}
						for ; w2 != cLength2; w2++ {
							if a3 = a2[w2]; a3 == nil {
								break loop
							}
							for ; w3 != cLength3; w3++ {
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
	result := (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + int32(bits.TrailingZeros(uint(nword)))
	if result == math.MaxInt32 {
		return -1
	}
	return result
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

func (bs *BitSet) NextSetBit(i int32) int32 {
	/*  The index value (i) of this method is permitted to be Integer.MAX_VALUE,
	as this is needed to make the loop defined above work: just in case the
	bit labelled Integer.MAX_VALUE-1 is set. This case is not optimised:
	but eventually -1 will be returned, as this will be included with
	any search that goes off the end of the level1 array. */

	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	/*  This is the word from which the search begins. */
	w := i >> cShift3
	w3 := w & cMask3
	w2 := (w >> cShift2) & cMask2
	w1 := w >> cShift1

	word := wordType(0)
	aLength := int32(len(bs.bits))

	/*  Is the next set bit in the same word at the nominated beginning bit
	(including the nominated beginning bit itself). The first check is
	whether the starting bit is within the structure at all. */
	var a2 b2DimType
	var a3 b1DimType
	if w1 < aLength {
		a2 = bs.bits[w1]
		result := true
		if a2 != nil {
			if a3 = a2[w2]; a3 != nil {
				word = a3[w3] & wordType(^uint(0)<<remainderOf64(i))
				result = word == 0
			}
		}
		if result {
			/*  So now start a search though the rest of the entries for a bit. */
			w++
			w3 = w & cMask3
			w2 = (w >> cShift2) & cMask2
			w1 = w >> cShift1
		major:
			for ; w1 != aLength; w1++ {
				if a2 = bs.bits[w1]; a2 != nil {
					for ; w2 != cLength2; w2++ {
						if a3 = a2[w2]; a3 != nil {
							for ; w3 != cLength3; w3++ {
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
	}
	return (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + int32(bits.TrailingZeros(uint(word)))

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
func (bs *BitSet) PreviousClearBit(i int32) int32 {
	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}

	bits := bs.bits
	aSize := int32(len(bs.bits) - 1)

	w := i >> cShift3
	w3 := w & cMask3
	w2 := (w >> cShift2) & cMask2
	w1 := w >> cShift1
	if w1 > aSize {
		return i
	}

	if aSize < w1 {
		w1 = aSize
	}
	w4 := i % cLength4

	word := wordType(0)
	var a2 b2DimType
	var a3 b1DimType

	for ; w1 >= 0; w1-- {
		if a2 = bits[w1]; a2 == nil {
			return (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + w4
		}
		for ; w2 >= 0; w2-- {
			if a3 = a2[w2]; a3 == nil {
				return (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + w4
			}
			for ; w3 >= 0; w3-- {
				if word = a3[w3]; word == 0 {
					return (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + w4
				}
				for bitIdx := w4; bitIdx >= 0; bitIdx-- {
					if t := word & (1 << uint(bitIdx)); t == 0 {
						return (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + bitIdx
					}
				}
				w4 = cLength4Size
			}
			w3 = cLength3Size
		}
		w2 = cLength2Size
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

func (bs *BitSet) PreviousSetBit(i int32) int32 {
	if i < 0 {
		panic(fmt.Sprintf("IndexOutOfBoundsException(i=%v)", i))
	}
	bits := bs.bits
	aSize := int32(len(bs.bits) - 1)

	/*  This is the word from which the search begins. */
	w := i >> cShift3
	w1 := w >> cShift1
	var w2, w3, w4 int32
	/*  But if its off the end of the array, start from the very end. */
	if w1 > aSize {
		w1 = aSize
		w2 = cLength2Size
		w3 = cLength3Size
		w4 = cLength4Size
	} else {
		w2 = (w >> cShift2) & cMask2
		w3 = w & cMask3
		w4 = i % cLength4
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
								if t := word & (1 << remainderOf64(bitIdx)); t != 0 {
									return (((w1 << cShift1) + (w2 << cShift2) + w3) << cShift3) + bitIdx
								}
							}
						}
						w4 = cLength4Size
					}
				}
				w3 = cLength3Size
				w4 = cLength4Size
			}
		}
		w2 = cLength2Size
		w3 = cLength3Size
		w4 = cLength4Size
	}
	return -1
}

func (bs *BitSet) ClearAll() {
	/*  This simply resets to null all the entries in the set. */
	bs.nullify(0)
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
func (bs BitSet) String() string {
	var p = "{"
	i := bs.NextSetBit(0)
	/*  Loop so long as there is another bit to append to the String. */
	for i >= 0 {
		/*  Append that next bit */
		p += fmt.Sprintf("%v", i)
		/*  Find the position of the next bit to show. */
		j := bs.NextSetBit(i + 1)
		if bs.compactionCount > 0 {
			/*  Give up if there is no next bit to show. */
			if j < 0 {
				break
			}
			/*  Find the next clear bit is after the current bit, i.e., i */
			last := bs.NextClearBit(i)
			/*  Compute the position of the next clear bit after the current
			subsequence of set bits. */
			if last < 0 {
				last = math.MaxInt32
			}
			/*  If the subsequence is more than the specified bits long, then
			collapse the subsequence into one entry in the String. */
			if (i + bs.compactionCount) < last {
				p += fmt.Sprintf("..%v", last-1)
				/*  Having accounted for a subsequence of bits that are all set,
				recompute the label of the next bit to show. */
				j = bs.NextSetBit(last)
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

/**
 *  Returns the number of bits set to <code>true</code> in this
 *  <code>SparseBitSet</code>.
 *
 * @return      the number of bits set to true in this SparseBitSet
 * @since       1.6
 */

//Cardinality ...
func (bs *BitSet) Cardinality() int32 {
	bs.statisticsUpdate() // Update size, cardinality and length values
	return bs.cache.cardinality
}

/**
 *  Convenience method for statistics if the individual results are not needed.
 *
 * @return      a String detailing the statistics of the bit set
 * @see         #statistics(String[])
 * @since       1.6
 */

//StatisticsAll ...
func (bs *BitSet) StatisticsAll() string {
	return bs.Statistics(nil)
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
func (bs *BitSet) Statistics(values []string) string {
	bs.statisticsUpdate() //  Ensure statistics are up-to-date
	v := make([]string, Statistics_Values_Length)
	/*  Assign the statistics values to the appropriate entry. The order
	of the assignments does not matter--the ordinal serves to get the
	values into the matching order with the labels from the enumeration. */
	v[Size] = strconv.Itoa(int(bs.Size()))
	v[Length] = strconv.Itoa(int(bs.Length()))
	v[Cardinality] = strconv.Itoa(int(bs.Cardinality()))
	v[Total_words] = strconv.Itoa(int(bs.cache.count))
	v[Set_array_length] = strconv.Itoa(len(bs.bits))
	v[Set_array_max_length] = strconv.Itoa(int(cMaxLength1))
	v[Level2_areas] = strconv.Itoa(int(bs.cache.a2Count))
	v[Level2_area_length] = strconv.Itoa(int(cLength2))
	v[Level3_blocks] = strconv.Itoa(int(bs.cache.a3Count))
	v[Level3_block_length] = strconv.Itoa(int(cLength3))
	v[Compaction_count_value] = strconv.Itoa(int(bs.compactionCount))

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
