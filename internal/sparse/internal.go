package sparse

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
)

func isZeroBlock(a3 b1DimType) bool {
	for _, word := range a3 {
		if word != 0 {
			return false
		}
	}
	return true
}

func remainderOf64(a int32) uint {
	return uint(a % 64)
}

func highestOneBit(x int32) (result int32) {
	l := bits.Len(uint(x))
	if l == 0 {
		return
	}
	result = 1 << (uint(l) - 1)
	return
}

func (bs *BitSet) resize(index int32) {
	/*  Find an array size that is a power of two that is as least as large
	enough to contain the index requested. */
	w1 := int32((index >> cShift3) >> cShift3)
	newSize := int32(highestOneBit(w1))
	if newSize == 0 {
		newSize = 1
	}
	if w1 >= newSize {
		newSize <<= 1
	}
	if newSize > cMaxLength1 {
		newSize = cMaxLength1
	}

	aLength1 := int32(0)
	if bs.bits != nil {
		aLength1 = int32(len(bs.bits))
	}

	if newSize != aLength1 || bs.bits == nil {
		// only if the size needs to be changed
		temp := make(b3DimType, newSize) //  Get the new array
		if aLength1 != 0 {
			/*  If it exists, copy old array to the new array. */
			copy(temp, bs.bits)
			bs.nullify(0) //  Don't leave unused pointers around. */
		}
		bs.bits = temp                //  Set new array as the set array
		bs.bitsLength = math.MaxInt32 //  Index of last possible bit, plus one.
		if newSize != cMaxLength1 {
			bs.bitsLength = newSize * cUnit
		}
	}
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
func (bs *BitSet) throwIndexOutOfBoundsException(i, j uint32) {
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

func (bs *BitSet) nullify(start int32) {
	aLength := int32(len(bs.bits))
	if start < aLength {
		for w := start; w != aLength; w++ {
			bs.bits[w] = nil
		}
		bs.cache.hash = 0 //  Invalidate size, etc., values
	}
}

/**
 *  Intializes all the additional objects required for correct operation.
 *
 * @since       1.6
 */
func (bs *BitSet) constructorHelper() {
	bs.spare = make(b1DimType, cLength3)
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

type strateger interface {
	properties() int32
	start(*BitSet) bool
	word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool
	block(base, u3, v3 int32, a3, b3 b1DimType) bool
	finish(cache *cacheType, a2Count, a3Count int32)
}

func (bs *BitSet) setScanner(i, j int32, b *BitSet, op strateger) {

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
		bs.cache.hash = 0
	}

	if j < i || (i+1) < 1 {
		panic(fmt.Sprintf("throwIndexOutOfBoundsException(%v,%v)", i, j))
	}

	if i == j {
		return
	}

	/*  Get the values of all the short-cut options. */
	properties := op.properties()
	falseOpFalseEqFalse := (properties & cFalseOpFalseEqFalse) != 0
	falseOpValueEqFalse := (properties & cFalseOpValueEqFalse) != 0
	valueOpFalseEqFalse := (properties & cValueOpFalseEqFalse) != 0
	valueOpFalseEqValue := (properties & cValueOpFalseEqValue) != 0

	/*  Index of the current word, and mask for the first word,
	to be processed in the bit set. */
	u := int32(i) >> cShift3
	//final long um = ~0L << i;
	um := wordType(^uint(0) << remainderOf64(i))

	/*  Index of the final word, and mask for the final word,
	to be processed in the bit set. */
	v := int32((j - 1)) >> cShift3
	// final long vm = ~0L >>> -j;
	vm := ^wordType(^uint(0) << remainderOf64(j))

	/*  Set up the two bit arrays (if the second exists), and their
	corresponding lengths (if any). */
	a1 := bs.bits //  Level1, i.e., the bit arrays
	aLength1 := int32(len(bs.bits))

	var b1 b3DimType
	var bLength1 int32

	if b != nil {
		b1 = b.bits
		bLength1 = int32(len(b.bits))
	}

	/*  Calculate the initial values of the parts of the words addresses,
	as well as the location of the final block to be processed.  */
	u1 := u >> cShift1
	u2 := (u >> cShift2) & cMask2
	u3 := u & cMask3
	v1 := v >> cShift1
	v2 := (v >> cShift2) & cMask2
	v3 := v & cMask3
	lastA3Block := (v1 << cLevel2) + v2

	/*  Initialize the local copies of the counts of blocks and areas; and
	whether there is a partial first block.  */
	var a2CountLocal int32
	var a3CountLocal int32
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
			(!haveA2 && !haveB2 && falseOpFalseEqFalse ||
				!haveA2 && falseOpValueEqFalse ||
				!haveB2 && valueOpFalseEqFalse) {
			//nested if!
			if u1 < aLength1 {
				a1[u1] = nil
			}
		} else {
			limit2 := cLength2
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

				a3Block := (u1 << cLevel2) + u2
				notLastBlock := lastA3Block != a3Block

				/*  Handling of level 3 empty areas: determined by the
				properties of the strategy. */
				if (!haveA3 && !haveB3 && falseOpFalseEqFalse || !haveA3 && falseOpValueEqFalse || !haveB3 && valueOpFalseEqFalse) && notFirstBlock && notLastBlock {
					/*  Do not need level3 block, so remove it, and move on. */
					if haveA2 {
						a2[u2] = nil
					}
				} else {
					/*  So what is needed is the level3 block. */
					base3 := a3Block << cShift2
					limit3 := cLength3
					if !notLastBlock {
						limit3 = int32(v3)
					}
					if !haveA3 {
						a3 = bs.spare
						a3IsSpare = true
					}
					if !haveB3 {
						b3 = iZeroBlock
					}
					isZero := false
					if notFirstBlock && notLastBlock {
						if valueOpFalseEqValue && !haveB3 {
							isZero = isZeroBlock(a3)
						} else {
							isZero = op.block(base3, 0, cLength3, a3, b3)
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
								if limit3 != cLength3 {
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
							if int32(i) >= bs.bitsLength { //Check that the set is large
								//  enough to take the new block
								bs.resize(i) //  Make it large enough
								a1 = bs.bits //  Update reference and length
								aLength1 = int32(len(a1))
							}
							if a2 == nil { //  Ensure a level 2 area
								a2 = make(b2DimType, cLength2)
								a1[u1] = a2
								haveA2 = true //  Ensure know level2 not empty
							}
							a2[u2] = a3 //  Insert the level3 block
							a3IsSpare = false
							bs.spare = make(b1DimType, cLength3) // Replace the spare

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
			if u2 == cLength2 && a2IsEmpty && u1 < aLength1 {
				a1[u1] = nil
			} else {
				a2CountLocal++ //  Count level 2 areas
			}
		}
		/*  Advance the value of u based on what happened. */
		u1++
		u = (u1 << cShift1)
		i = u << cShift3
		u2 = 0 //  u3 = 0
		//  Compute next word and bit index
		if i < 0 {
			i = math.MaxInt32 //  Don't go over the end
		}

	} /* end while( i < j ) */

	/*  Do whatever the strategy needs in order to finish. */
	op.finish(&bs.cache, a2CountLocal, a3CountLocal)
}

/**
 *  The entirety of the bit set is examined, and the various statistics of
 *  the bit set (size, length, cardinality, hashCode, etc.) are computed. Level
 *  arrays that are empty (i.e., all zero at level 3, all null at level 2) are
 *  replaced by null references, ensuring a normalized representation.
 *
 * @since       1.6
 */
func (bs *BitSet) statisticsUpdate() {
	if bs.cache.hash != 0 {
		return
	}
	bs.setScanner(0, bs.bitsLength, nil, new(updateStrategyType))
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
func (bs *BitSet) clone() (result *BitSet) {
	result = new(BitSet)

	reflect.Copy(reflect.ValueOf(*result), reflect.ValueOf(*bs))

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

	result.setScanner(0, bs.bitsLength, bs, copyStrategy)
	return result
}
