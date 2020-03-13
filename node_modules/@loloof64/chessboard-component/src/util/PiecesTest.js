export function cellAlgebraic(file, rank) {
    return `${String.fromCharCode('a'.charCodeAt(0) + file)}${String.fromCharCode('1'.charCodeAt(0) + rank)}`;
}

export function getPieceAt(chessLogic, file, rank) {
    const pieceValue = chessLogic.get(cellAlgebraic(file, rank));
    return pieceValue;
}

export function isWhitePawnAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'p' && piece.color === 'w';
}

export function isWhiteKnightAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'n' && piece.color === 'w';
}

export function isWhiteBishopAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'b' && piece.color === 'w';
}

export function isWhiteRookAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'r' && piece.color === 'w';
}

export function isWhiteQueenAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'q' && piece.color === 'w';
}

export function isWhiteKingAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'k' && piece.color === 'w';
}

export function isBlackPawnAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'p' && piece.color === 'b';
}

export function isBlackKnightAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'n' && piece.color === 'b';
}

export function isBlackBishopAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'b' && piece.color === 'b';
}

export function isBlackRookAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'r' && piece.color === 'b';
}

export function isBlackQueenAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'q' && piece.color === 'b';
}

export function isBlackKingAtCell(chessLogic, file, rank) {
    const piece = getPieceAt(chessLogic, file, rank);
    if ([null, undefined].includes(piece)) return false;
    return piece.type === 'k' && piece.color === 'b';
}
