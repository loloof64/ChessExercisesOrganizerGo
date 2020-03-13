import {getPieceAt, cellAlgebraic} from './PiecesTest';
import { Chess } from 'chess.js';

function getCell({x, y, cellsSize, reversed}) {
    const cellX = Math.floor((x - cellsSize * 0.5) / cellsSize);
    const cellY = Math.floor((y - cellsSize * 0.5) / cellsSize);

    const file = [true, "true"].includes(reversed) ? 7-cellX : cellX;
    const rank = [true, "true"].includes(reversed)  ? cellY : 7-cellY;

    return [file, rank];
}

function getLocalCoordinates(event, rootElement) {
    if (!rootElement) return;

    const thisComponentLocation = rootElement.getBoundingClientRect();
    const localX = event.clientX - thisComponentLocation.x;
    const localY = event.clientY - thisComponentLocation.y;

    return [localX, localY];
}

export function handleMouseDown({event, cellsSize, reversed, rootElement,
     logic, setupDnd, gameInProgress, waitingForExternalMove}) {
    if (!gameInProgress) return;
    if (waitingForExternalMove) return;
    const [x, y] = getLocalCoordinates(event, rootElement);
    const [file, rank] = getCell({x,y,cellsSize, reversed});

    const piece = getPieceAt(logic, file, rank);
    if ([null, undefined].includes(piece)) return;

    setupDnd({x, y, file, rank, piece});
}

export function handleMouseMove({event, dragAndDropInProgress, updateDndLocation, 
    rootElement, cancelDnd, cellsSize, reversed, promotionPending,
     gameInProgress, waitingForExternalMove}) {
    if (!gameInProgress) return;
    if (waitingForExternalMove) return;
    if (!dragAndDropInProgress) return;
    if (promotionPending) return;

    const [x, y] = getLocalCoordinates(event, rootElement);
    const [file, rank] = getCell({x,y,cellsSize, reversed});

    const thisComponentLocation = rootElement.getBoundingClientRect();
    const inBounds = x >= 0 && x <= thisComponentLocation.width &&
        y >= 0 && y <= thisComponentLocation.height;

    if (!inBounds) {
        cancelDnd();
        return;
    }

    updateDndLocation(x,y, file, rank);
}

export function handleMouseUp({event, cellsSize, reversed, rootElement, logic,
     dragAndDropInProgress, dndPieceData, cancelDnd, updateLogic,
     updateAndEmitLastMove, promotionPending, setPromotionPending, gameInProgress,
     handleGameEndedStatus, updateWaitingForExternalMove, waitingForExternalMove}) {
    if (!gameInProgress) return;
    if (waitingForExternalMove) return;
    if (!dragAndDropInProgress) return;
    if (promotionPending) return;
    
    const [x, y] = getLocalCoordinates(event, rootElement);
    const [file, rank] = getCell({x,y,cellsSize, reversed});

    const originCell = dndPieceData.originCell;

    const inBounds = originCell.file >= 0 && originCell.file <= 7 &&
        originCell.rank >= 0 && originCell.rank <= 7 &&
        file >= 0 && file <= 7 && rank >= 0 && rank <= 7;

    if (! inBounds) {
        cancelDnd();
        return;
    }

    const testMoveObject = buildMoveObject({startFile: originCell.file, startRank: originCell.rank,
        endFile: file, endRank: rank});
    const logicClone = new Chess(logic.fen());
    const promotionTest = logicClone.move(testMoveObject);

    const validTestMove = ![null, undefined].includes(promotionTest);
    const pieceAtOriginCell = logic.get(cellAlgebraic(originCell.file, originCell.rank));
    const isPawnMoving = pieceAtOriginCell && pieceAtOriginCell.type === 'p';
    const targetOnFirstOrLastRank = rank === 0 || rank === 7;

    const isPromotionMove = validTestMove && isPawnMoving && targetOnFirstOrLastRank;

    if (isPromotionMove) {
        setPromotionPending({
            startFile: originCell.file, startRank: originCell.rank,
            endFile: file, endRank: rank
        });
        return;
    }

    const moveObject = buildMoveObject({startFile: originCell.file, startRank: originCell.rank,
         endFile: file, endRank: rank});

    const logicBeforeMove = new Chess(logic.fen());
    const result = logic.move(moveObject);

    if ([null, undefined].includes(result)) {
        cancelDnd();
        return;
    }

    updateLogic();
    updateAndEmitLastMove({startFile: originCell.file, startRank: originCell.rank,
        endFile: file, endRank: rank, logicBeforeMove, logicAfterMove: logic});
    cancelDnd();

    handleGameEndedStatus();
    updateWaitingForExternalMove();
}

export function handleMouseExited({event, cancelDnd, promotionPending, 
    gameInProgress, waitingForExternalMove}) {
    if (!gameInProgress) return;
    if (waitingForExternalMove) return;
    if (promotionPending) return;
    cancelDnd();
}

function buildMoveObject({startFile, startRank, endFile, endRank, promotion = 'q'}) {
    const startAlgebraic = cellAlgebraic(startFile, startRank);
    const endAlgebraic = cellAlgebraic(endFile, endRank);

    return {
        from: startAlgebraic,
        to: endAlgebraic,
        promotion,
    };
}


export function isDnDOriginCell(dndPieceData, file, rank) {
    if ([undefined, null].includes(dndPieceData)) return undefined;
    
    const originCell = dndPieceData.originCell;
    if ([undefined, null].includes(originCell)) return undefined;

    return originCell.file === file && originCell.rank === rank;
}

export function isWhitePawnDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'p' && dndPieceData.color === 'w';
}

export function isWhiteKnightDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'n' && dndPieceData.color === 'w';
}

export function isWhiteBishopDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'b' && dndPieceData.color === 'w';
}

export function isWhiteRookDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'r' && dndPieceData.color === 'w';
}

export function isWhiteQueenDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'q' && dndPieceData.color === 'w';
}

export function isWhiteKingDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'k' && dndPieceData.color === 'w';
}

export function isBlackPawnDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'p' && dndPieceData.color === 'b';
}

export function isBlackKnightDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'n' && dndPieceData.color === 'b';
}

export function isBlackBishopDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'b' && dndPieceData.color === 'b';
}

export function isBlackRookDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'r' && dndPieceData.color === 'b';
}

export function isBlackQueenDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'q' && dndPieceData.color === 'b';
}

export function isBlackKingDragged(dndPieceData) {
    if ([undefined, null].includes(dndPieceData)) return false;
    return dndPieceData.type === 'k' && dndPieceData.color === 'b';
}