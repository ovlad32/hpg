package sparse

import (
	"fmt"
	"math"
	"math/bits"
)

type b3DimType [][][]int64
type b2DimType [][]int64
type b1DimType []int64

const compactionCountDefault uint32 = 2

/**
 *  The number of bits in a long value.
 */
const LENGTH4 uint32 = 64 //Long.SIZE

/**
 *  The number of bits in a positive integer, and the size of permitted index
 *  of a bit in the bit set.
 */
const INDEX_SIZE uint32 = 31 //Integer.SIZE - 1;

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
const MAX_LENGTH1 uint32 = 1 << LEVEL1

/**
 *  LENGTH2 is the number of entries in the any level2 area.
 */
const LENGTH2 uint32 = 1 << LEVEL2

/**
 *  LENGTH3 is the number of entries in the any level3 block.
 */
const LENGTH3 uint32 = 1 << LEVEL3

/**
 *  The shift to create the word index. (I.e., move it to the right end)
 */
const SHIFT3 uint32 = LEVEL4

/**
 *  MASK3 is the mask to extract the LEVEL3 address from a word index
 *  (after shifting by SHIFT3).
 */
const MASK3 uint32 = LENGTH3 - 1

/**
 *  SHIFT2 is the shift to bring the level2 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const SHIFT2 uint32 = LEVEL3

/**
 *  UNIT is the greatest number of bits that can be held in one level1 entry.
 *  That is, bits per word by words per level3 block by blocks per level2 area.
 */
var UNIT uint32 = LENGTH2 * LENGTH3 * LENGTH4

/**
 *  MASK2 is the mask to extract the LEVEL2 address from a word index
 *  (after shifting by SHIFT3 and SHIFT2).
 */
const MASK2 uint32 = LENGTH2 - 1

/**
 *  SHIFT1 is the shift to bring the level1 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const SHIFT1 uint32 = LEVEL2 + LEVEL3

/**
 *  LENGTH2_SIZE is maximum index of a LEVEL2 page.
 */
const LENGTH2_SIZE uint32 = LENGTH2 - 1

/**
 *  LENGTH3_SIZE is maximum index of a LEVEL3 page.
 */
const LENGTH3_SIZE uint32 = LENGTH3 - 1

/**
 *  LENGTH4_SIZE is maximum index of a bit in a LEVEL4 word.
 */
const LENGTH4_SIZE uint32 = LENGTH4 - 1

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
	compactionCount uint32

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
	bitsLength uint32
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

func New() *BitSet {
	return NewWithSizeAndCompactionCount(1, compactionCountDefault)
}
func NewWithSize(capacity uint32) *BitSet {
	return NewWithSizeAndCompactionCount(1, compactionCountDefault)
}

func NewWithSizeAndCompactionCount(capacity uint32, compactionCount uint32) *BitSet {
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

func highestOneBit(x uint32) (result uint32) {
	l := bits.Len32(x)
	if result == 0 {
		return
	}
	result = 1 << (uint(l) - 1)
	return
}

func (this *BitSet) resize(index uint32) {
	/*  Find an array size that is a power of two that is as least as large
	enough to contain the index requested. */
	w1 := (index >> SHIFT3) >> SHIFT1
	newSize := highestOneBit(w1)
	if newSize == 0 {
		newSize = 1
	}
	if w1 >= newSize {
		newSize <<= 1
	}
	if newSize > MAX_LENGTH1 {
		newSize = MAX_LENGTH1
	}

	aLength1 := uint32(0)
	if this.bits != nil {
		aLength1 = uint32(len(this.bits))
	}

	if newSize != aLength1 || this.bits == nil {
		// only if the size needs to be changed
		temp := make(b3DimType, newSize) //  Get the new array
		if aLength1 != 0 {
			/*  If it exists, copy old array to the new array. */
			copy(temp, this.bits)
			this.nullify(0) //  Don't leave unused pointers around. */
		}
		this.bits = temp                        //  Set new array as the set array
		this.bitsLength = uint32(math.MaxInt32) //  Index of last possible bit, plus one.
		if newSize != MAX_LENGTH1 {
			this.bitsLength = newSize * UNIT
		}
	}
}

func (this *BitSet) nullify(start int) {
	aLength := len(this.bits)
	if start < aLength {
		for w := start; w != aLength; w++ {
			this.bits[w] = nil
		}
		this.cache.hash = 0 //  Invalidate size, etc., values
	}
}

func (this *BitSet) constructorHelper() {
	this.spare = make(b1DimType, LENGTH3)
	this.cache = new(cacheType)
	this.updateStrategy = new(updateStrategyType)
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
	u := i >> SHIFT3
	um := ^0 << uint32(i)

	/*  Index of the final word, and mask for the final word,
	to be processed in the bit set. */
	v := (j - 1) >> SHIFT3
	vm := bits.RotateLeft32(int32(^0), j)

	/*  Set up the two bit arrays (if the second exists), and their
	corresponding lengths (if any). */
	a1 := this.bits //  Level1, i.e., the bit arrays
	aLength1 := this.bits.length

	var b1 b3DimType = nil
	bLength1 := 0

	if b1 != nil {
		b1 = b.bits
		bLength1 = len(b.bits)
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
	a2CountLocal = 0
	a3CountLocal = 0
	notFirstBlock = u == 0 && um == ^0

	/*  The first level2 is cannot be judged empty if not being scanned from
	the beginning. */
	a2IsEmpty := u2 == 0 //  Presumption
	for i < j {
		/*  Determine if there is a level2 area in both the a and the b set,
		and if so, set the references to these areas. */
		a2 = a1[u1]
		haveA2 := u1 < aLength1 && a2 != nil
		b2 = b1[u1]

		haveB2 = u1 < bLength1 && b1 != null && b2 != nil
		/*  Handling of level 2 empty areas: determined by the
		properties of the strategy. It is necessary to actually visit
		the first and last blocks of a scan, since not all of the block
		might participate in the operation, hence making decision based
		on just the references to the blocks could be wrong. */
		if (!haveA2 && !haveB2 && f_op_f_eq_f || !haveA2 && f_op_x_eq_f || !haveB2 && x_op_f_eq_f) && notFirstBlock && u1 != v1 {
			//nested if!
			if u1 < aLength1 {
				a1[u1] = nil
			}
		} else {
			limit2 := LENGTH2
			if u1 == v1 {
				limit2 = v2 + 1
			}

			for u2 != limit2 {
				/*  Similar logic applied here as for the level2 blocks.
				The initial and final block must be examined. In other
				cases, it may be possible to make a decision based on
				the value of the references, as indicated by the
				properties of the strategy. */
				a3 := a2[u2]
				haveA3 = haveA2 && a3 != nil
				b3 := b2[u2]
				haveB3 = haveB2 && b3 != nil
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
						limit3 = v3
					}
					if !haveA3 {
						a3 = spare
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
								isZero = op.word(base3, u3, a3, b3, um&vm)
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
						if a3 == spare {
							if i >= bitsLength { //Check that the set is large
								//  enough to take the new block
								resize(i) //  Make it large enough
								a1 = bits //  Update reference and length
								aLength1 = a1.length
							}
							if a2 == null { //  Ensure a level 2 area
								a2 = make(b2DimType, LENGTH2)
								a1[u1] = a2
								haveA2 = true //  Ensure know level2 not empty
							}
							a2[u2] = a3                      //  Insert the level3 block
							spare = make(b1DimType, LENGTH3) // Replace the spare
						}
						a3CountLocal++ // Count the level 3 block
					}
					a2IsEmpty &= !(haveA2 && a2[u2] != null)
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
			i = Math.MaxInt32 //  Don't go over the end
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
func (b *BitSet) statisticsUpdate() {
	if b.cache.hash != 0 {
		return
	}
	b.setScanner(0, b.bitsLength, nil, b.updateStrategy)
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
type statistics int

const (
	/**
	 *  The size of the bit set, as give by the <i>size</i>() method.
	 */
	Size statistics = iota // 0
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
)

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
	hash uint

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
	word(base, u3 int32, a3, b3 []int64, mask int64) bool
	block(base, u3, v3 int32, a3, b3 []int64) bool
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

func isZeroBlock(a3 []int64) bool {
	for _, word := range a3 {
	}
	if word != 0 {
		return false
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

func (st andStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3]&b3[u3] | ^mask
	return a3[u3] == 0
}

func (st andStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
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

func (st andNotStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3] & ^(b3[u3] & mask)
	return a3[u3] == 0
}

func (st andNotStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
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

func (st clearStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3] & ^mask
	return a3[u3] == 0
}

func (st clearStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
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

func (st copyStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = b3[u3] & mask
	return a3[u3] == 0
}

func (st copyStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
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

func (st *equalsStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	word := a3[u3]
	st.result = st.result && ((word & mask) == (b3[u3] & mask))
	return word == 0
}

func (st *equalsStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
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

func (st flipStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3] ^ mask
	return a3[u3] == 0
}

func (st flipStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] ^ (^0)
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

func (st *intersectsStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	word := a3[u3]
	st.result = st.result || ((word & b3[u3] & mask) != 0)
	return word == 0
}

func (st *intersectsStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		word := a3[w3]
		st.result = st.result || (word & b3[w3])
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

func (st orStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3] | (b3[w3] & mask)
	return a3[u3] == 0
}

func (st orStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
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

func (st setStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3] | mask
	return a3[u3] == 0
}

func (st setStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = ^0
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
	wordMin int64

	/**
	 *  Working space for find the size and length of the bit set. Holds the
	 *  index of the last non-empty word in the set.
	 */
	wMax int32

	/**
	 *  Working space for find the size and length of the bit set. Holds a copy
	 *  of the last non-empty word in the set.
	 */
	wordMax int64

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
	count int

	/**
	 *  Working space for counting the number of non-zero bits in the bit set.
	 *  Holds the current state of the computation of the cardinality.This
	 *  value is ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	cardinality int
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

func (st updateStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	word := a3[u3]
	word1 := word & mask
	if word1 != 0 {
		compute(base+u3, word1)
	}
	return word == 0
}

func (st *updateStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
	isZero = true //  Presumption
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		if word := a3[w3]; word != 0 {
			isZero = false
			compute(base+w3, word)
		}
	}
	isZero = false
	return
}

func (st *updateStrategyType) finish(a2Count, a3Count int32) {
	cache.a2Count = a2Count
	cache.a3Count = a3Count
	cache.count = st.count
	cache.cardinality = st.cardinality
	cache.length = (st.wMax+1)*LENGTH4 - bits.LeadingZeros(st.wordMax)
	cache.size = cache.length - st.wMin*LENGTH4 - bits.LeadingZeros(st.wordMin)
	cache.hash = (int)((st.hash >> Integer.SIZE) ^ st.hash)
}

func (st *updateStrategyType) compute(index int32, word int64) {
	/*  Count the number of actual words being used. */
	count++
	/*  Continue to accumulate the hash value of the set. */
	st.hash = st.hash ^ (word * (long)(index+1))
	/*  The first non-zero word contains the first actual bit of the
	    set. The location of this bit is used to compute the set size. */
	if wMin < 0 {
		st.wMin = index
		st.wordMin = word
	}
	/*  The last non-zero word contains the last actual bit of the set.
	    The location of this bit is used to compute the set length. */
	st.wMax = index
	st.wordMax = word
	/*  Count the actual bits, so as to get the cardinality of the set. */
	st.cardinality = st.cardinality + bits.OnesCount(word)
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

func (st xorStrategyType) word(base, u3 int32, a3, b3 []int64, mask int64) bool {
	a3[u3] = a3[u3] ^ (b3[u3] & mask)
	return a3[u3] == 0
}

func (st xorStrategyType) block(base, u3, v3 int32, a3, b3 []int64) (isZero bool) {
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] ^ b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	isZero = false
	return
}
func (st xorStrategyType) finish(a2Count, a3Count int32) {}
