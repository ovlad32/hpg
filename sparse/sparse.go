package sparse

import (
	"math"
	"math/bits"
)

const compactionCountDefault int = 2

/**
 *  The number of bits in a long value.
 */
const LENGTH4 int64 = 64 //Long.SIZE

/**
 *  The number of bits in a positive integer, and the size of permitted index
 *  of a bit in the bit set.
 */
const INDEX_SIZE = 31 //Integer.SIZE - 1;

/**
 *  The label (index) of a bit in the bit set is essentially broken into
 *  4 "levels". Respectively (from the least significant end), level4, the
 *  address within word, the address within a level3 block, the address within
 *  a level2 area, and the level1 address of that area within the set.
 *
 *  LEVEL4 is the number of bits of the level4 address (number of bits need
 *  to address the bits in a long)
 */
const LEVEL4 = 6

/**
 *  LEVEL3 is the number of bits of the level3 address.
 *  Do not change!
 */
const LEVEL3 = 5

/**
 *  LEVEL2 is the number of bits of the level3 address.
 *  Do not change!
 */
const LEVEL2 = 5

/**
 *  LEVEL1 is the number of bits of the level1 address.
 */
const LEVEL1 = INDEX_SIZE - LEVEL2 - LEVEL3 - LEVEL4

/**
 *  MAX_LENGTH1 is the maximum number of entries in the level1 set array.
 */
const MAX_LENGTH1 = 1 << LEVEL1

/**
 *  LENGTH2 is the number of entries in the any level2 area.
 */
const LENGTH2 = 1 << LEVEL2

/**
 *  LENGTH3 is the number of entries in the any level3 block.
 */
const LENGTH3 = 1 << LEVEL3

/**
 *  The shift to create the word index. (I.e., move it to the right end)
 */
const SHIFT3 = LEVEL4

/**
 *  MASK3 is the mask to extract the LEVEL3 address from a word index
 *  (after shifting by SHIFT3).
 */
const MASK3 = LENGTH3 - 1

/**
 *  SHIFT2 is the shift to bring the level2 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const SHIFT2 = LEVEL3

/**
 *  UNIT is the greatest number of bits that can be held in one level1 entry.
 *  That is, bits per word by words per level3 block by blocks per level2 area.
 */
const UNIT = LENGTH2 * LENGTH3 * LENGTH4

/**
 *  MASK2 is the mask to extract the LEVEL2 address from a word index
 *  (after shifting by SHIFT3 and SHIFT2).
 */
const MASK2 = LENGTH2 - 1

/**
 *  SHIFT1 is the shift to bring the level1 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const SHIFT1 = LEVEL2 + LEVEL3

/**
 *  LENGTH2_SIZE is maximum index of a LEVEL2 page.
 */
const LENGTH2_SIZE = LENGTH2 - 1

/**
 *  LENGTH3_SIZE is maximum index of a LEVEL3 page.
 */
const LENGTH3_SIZE = LENGTH3 - 1

/**
 *  LENGTH4_SIZE is maximum index of a bit in a LEVEL4 word.
 */
const LENGTH4_SIZE = LENGTH4 - 1

/** An empty level 3 block is kept for use when scanning. When a source block
 *  is needed, and there is not already one in the corresponding bit set, the
 *  ZERO_BLOCK is used (as a read-only block). It is a source of zero values
 *  so that code does not have to test for a null level3 block. This is a
 *  static block shared everywhere.
 */
var ZERO_BLOCK = make([]int64, LENGTH3)


 /**
     *  Word and block <b>and</b> strategy.
     */
	 const andStrategy := newAndStrategy();
	 /**
	  *  Word and block <b>andNot</b> strategy.
	  */
	 const  andNotStrategy = newAndNotStrategy();
	 /**
	  *  Word and block <b>clear</b> strategy.
	  */
	 const clearStrategy = newClearStrategy();
	 /**
	  *  Word and block <b>copy</b> strategy.
	  */
	 const copyStrategy = newCopyStrategy();
	 /**
	  *  Word and block <b>flip</b> strategy.
	  */
	 const  flipStrategy = newFlipStrategy();
	 /**
	  *  Word and block <b>intersects</b> strategy.
	  */
	 const intersectsStrategy = newIntersectsStrategy();
	 /**
	  *  Word and block <b>or</b> strategy.
	  */
	  const orStrategy = newOrStrategy();
	 /**
	  *  Word and block <b>set</b> strategy.
	  */
	  const setStrategy = newSetStrategy();

	  /**
	  *  Word and block <b>xor</b> strategy.
	  */
	  const xorStrategy = newXorStrategy();

	 

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
	bits [][][]int64

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
	cache *Cache
	/**
	 *  A spare level 3 block is kept for use when scanning. When a target block
	 *  is needed, and there is not already one in the bit set, the spare is
	 *  provided. If non-zero values are placed into this block, it is moved to the
	 *  resulting set, and a new spare is acquired. Note: a new spare needs to
	 *  be allocated when the set is cloned (so that the spare is not shared
	 *  between two sets).
	 */
	spare []int64

	 /**
	  *  Word and block <b>equals</b> strategy.
	  */
	  equalsStrategy  *EqualsStrategy;
	  /**
	  *  Word and block <b>update</b> strategy.
	  */
	  updateStrategy *updateStrategy;
	 
}

func New() *BitSet {
	return NewWithSizeAndCompactionCount(1, compactionCountDefault)
}
func NewWithSize(capacity int32) *BitSet {
	return NewWithSizeAndCompactionCount(1, compactionCountDefault)
}

func NewWithSizeAndCompactionCount(capacity uint32, compactionCount uint32) *BitSet {
	result = &BitSet{
		compactionCount: compactionCount,
	}
	result.resize(capacity - 1) //  Resize takes last usable index
	this.compactionCount = compactionCount
	/*  Ensure there is a spare level 3 block for the use of the set scanner.*/
	result.constructorHelper()
	result.statisticsUpdate()
	return result
}

func higestOneBit(x uint) (result uint) {
	result = bits.Len(x)
	if result == 0 {
		return
	}
	result = 1 << (result - 1)
	return
}

func (b *BitSet) resize(index uint) {
	/*  Find an array size that is a power of two that is as least as large
	enough to contain the index requested. */
	const w1 = (index >> SHIFT3) >> SHIFT1
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

	aLength1 := 0
	if b.bits != nil {
		aLength1 := len(bits)
	}

	if newSize != aLength1 || bits == nil {
		// only if the size needs to be changed
		temp = make([][][]int, newSize) //  Get the new array
		if aLength1 != 0 {
			/*  If it exists, copy old array to the new array. */
			copy(temp, bits)
			b.nullify(0) //  Don't leave unused pointers around. */
		}
		b.bits = temp                  //  Set new array as the set array
		b.bitsLength = math.MaxInt32() //  Index of last possible bit, plus one.
		if newSize != MAX_LENGTH1 {
			b.bitsLength = newSize * UNIT
		}
	}
}

func (b *BitSet) constructorHelper() {
	b.spare = make([]int, LENGTH3)
	b.cache = newCache()
	b.updateStrategy = newUpdateStrategy()
}
func (b *BitSet) statisticsUpdate() {
	if b.cache.hash != 0 {
		return
	}
	b.setScanner(0, b.bitsLength, nil, b.updateStrategy)
}

type strateger  interface {
	properties() int32
	start(*BitSet) bool
	word(base, u3 int32,  a3, b3 []int64, mask int64) bool
	block(base, u3, v3 int32,  a3, b3 []int64)
}	

type strategy  interface {

type xorStrategy struct

protected static class XorStrategy extends AbstractStrategy
    {
        @Override
        //  XorStrategy
        protected int properties()
        {
            return F_OP_F_EQ_F + X_OP_F_EQ_X;
        }

        @Override
        //  XorStrategy
        protected boolean start(SparseBitSet b)
        {
            if (b == null)
                throw new NullPointerException();
            return true;
        }

        @Override
        protected boolean word(int base, int u3, long[] a3, long[] b3, long mask)
        {
            return (a3[u3] ^= b3[u3] & mask) == 0;
        }

        @Override
        //  XorStrategy
        protected boolean block(int base, int u3, int v3, long[] a3, long[] b3)
        {
            boolean isZero = true; //  Presumption
            for (int w3 = u3; w3 != v3; ++w3)
                isZero &= (a3[w3] ^= b3[w3]) == 0;
            return isZero;

        }
    }