package sparse

const compactionCountDefault int32 = 2

/**
 *  The number of bits in a long value.
 */
const cLength4 int32 = 64 //Long.SIZE
const cIntegerSize = 32   //Integer.SIZE

/**
 *  The number of bits in a positive integer, and the size of permitted index
 *  of a bit in the bit set.
 */
const cIndexSize uint32 = cIntegerSize - 1 //Integer.SIZE - 1;

/**
 *  The label (index) of a bit in the bit set is essentially broken into
 *  4 "levels". Respectively (from the least significant end), level4, the
 *  address within word, the address within a level3 block, the address within
 *  a level2 area, and the level1 address of that area within the set.
 *
 *  LEVEL4 is the number of bits of the level4 address (number of bits need
 *  to address the bits in a long)
 */
const cLevel4 uint32 = 6

/**
 *  LEVEL3 is the number of bits of the level3 address.
 *  Do not change!
 */
const cLevel3 uint32 = 5

/**
 *  LEVEL2 is the number of bits of the level3 address.
 *  Do not change!
 */
const cLevel2 uint32 = 5

/**
 *  LEVEL1 is the number of bits of the level1 address.
 */
const cLevel1 = cIndexSize - cLevel2 - cLevel3 - cLevel4

/**
 *  MAX_LENGTH1 is the maximum number of entries in the level1 set array.
 */
const cMaxLength1 int32 = 1 << uint(cLevel1)

/**
 *  LENGTH2 is the number of entries in the any level2 area.
 */
const cLength2 int32 = 1 << uint(cLevel2)

/**
 *  LENGTH3 is the number of entries in the any level3 block.
 */
const cLength3 int32 = 1 << uint(cLevel3)

/**
 *  The shift to create the word index. (I.e., move it to the right end)
 */
const cShift3 uint32 = cLevel4

/**
 *  MASK3 is the mask to extract the LEVEL3 address from a word index
 *  (after shifting by SHIFT3).
 */
const cMask3 int32 = int32(cLength3 - 1)

/**
 *  SHIFT2 is the shift to bring the level2 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const cShift2 uint32 = cLevel3

/**
 *  UNIT is the greatest number of bits that can be held in one level1 entry.
 *  That is, bits per word by words per level3 block by blocks per level2 area.
 */
const cUnit int32 = cLength2 * cLength3 * cLength4

/**
 *  MASK2 is the mask to extract the LEVEL2 address from a word index
 *  (after shifting by SHIFT3 and SHIFT2).
 */
const cMask2 int32 = int32(cLength2 - 1)

/**
 *  SHIFT1 is the shift to bring the level1 address (from the word index) to
 *  the right end (i.e., after shifting by SHIFT3).
 */
const cShift1 uint32 = cLevel2 + cLevel3

/**
 *  LENGTH2_SIZE is maximum index of a LEVEL2 page.
 */
const cLength2Size int32 = int32(cLength2 - 1)

/**
 *  LENGTH3_SIZE is maximum index of a LEVEL3 page.
 */
const cLength3Size int32 = int32(cLength3 - 1)

/**
 *  LENGTH4_SIZE is maximum index of a bit in a LEVEL4 word.
 */
const cLength4Size int32 = int32(cLength4 - 1)

/** If the operation requires that when matching level2 areas or level3
 *  blocks are null, that no action is required, then this property is
 *  required. Corresponds to the top-left entry in the logic diagram for the
 *  operation being 0. For all the defined actual logic operations ('and',
 *  'andNot', 'or', and 'xor', this will be true, because for all these,
 *  "false" op "false" = "false".
 */
const cFalseOpFalseEqFalse = 0x1

/** If when level2 areas or level3 areas from the this set are null will
 *  require that area or block to remain null, irrespective of the value of
 *  the matching structure from the other set, then this property is required.
 *  Corresponds to the first row in the logic diagram being all zeros. For
 *  example, this is true for 'and' as well as 'andNot', and for 'clear', since
 *  false" & "x" = "false", and "false" &! "x" = "false".
 */
const cFalseOpValueEqFalse = 0x2

/** If when level2 areas or level3 areas from the other set are null will
 *  require the matching area or block in this set to be set to null,
 *  irrespective of the current values in the matching structure from the
 *  this, then this property is required. Corresponds to the first column
 *  in the logic diagram being all zero. For example, this is true for
 *  'and', since "x" & "false" = "false", as well as for 'clear'.
 */
const cValueOpFalseEqFalse = 0x4

/** If when a level3 area from the other set is null will require the
 *  matching area or block in this set to be left as it is, then this property
 *  is required. Corresponds to the first column of the logic diagram being
 *  equal to the left hand operand column. For example, this is true for 'or',
 *  'xor', and 'andNot', since for all of these "x" op "false" = "x".
 */
const cValueOpFalseEqValue = 0x8
