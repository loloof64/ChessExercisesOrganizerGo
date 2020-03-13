<script>
import { createEventDispatcher, onMount } from 'svelte';
import { get_current_component } from "svelte/internal"

const component = get_current_component()
const svelteDispatch = createEventDispatcher()
const dispatch = (name, detail) => {
	svelteDispatch(name, detail)
	component.dispatchEvent && component.dispatchEvent(new CustomEvent(name, { detail }))
}

let rootElement;
let gameInProgress = false;
let waitingForExternalMove = false;

export let size = 100;
export let background = '#124589';
export let reversed = false;

export let white_cell_color = 'GoldenRod';
export let black_cell_color = 'brown';
export let coordinates_color = 'DarkOrange';

export let origin_cell_color = 'crimson';
export let target_cell_color = 'ForestGreen';
export let dnd_cross_color = 'DimGrey';

export let move_highlight_color = 'CadetBlue';

export let promotion_dialog_title = 'Select the promotion piece';

export let white_player_human = 'true';
export let black_player_human = 'true';

import WhitePawn from './pieces/WhitePawn.svelte';
import WhiteKnight from './pieces/WhiteKnight.svelte';
import WhiteBishop from './pieces/WhiteBishop.svelte';
import WhiteRook from './pieces/WhiteRook.svelte';
import WhiteQueen from './pieces/WhiteQueen.svelte';
import WhiteKing from './pieces/WhiteKing.svelte';

import BlackPawn from './pieces/BlackPawn.svelte';
import BlackKnight from './pieces/BlackKnight.svelte';
import BlackBishop from './pieces/BlackBishop.svelte';
import BlackRook from './pieces/BlackRook.svelte';
import BlackQueen from './pieces/BlackQueen.svelte';
import BlackKing from './pieces/BlackKing.svelte';

import {
    cellAlgebraic,
    isWhitePawnAtCell,
    isWhiteKnightAtCell,
    isWhiteBishopAtCell,
    isWhiteRookAtCell,
    isWhiteQueenAtCell,
    isWhiteKingAtCell,
    isBlackPawnAtCell,
    isBlackKnightAtCell,
    isBlackBishopAtCell,
    isBlackRookAtCell,
    isBlackQueenAtCell,
    isBlackKingAtCell,
} from './util/PiecesTest.js';

import {
    handleMouseDown,
    handleMouseExited,
    handleMouseMove,
    handleMouseUp,
    isDnDOriginCell,
    isWhitePawnDragged,
    isWhiteKnightDragged,
    isWhiteBishopDragged,
    isWhiteRookDragged,
    isWhiteQueenDragged,
    isWhiteKingDragged,
    isBlackPawnDragged,
    isBlackKnightDragged,
    isBlackBishopDragged,
    isBlackRookDragged,
    isBlackQueenDragged,
    isBlackKingDragged,
} from './util/DragAndDrop.js';

import {Chess} from 'chess.js';


let logic = new Chess('8/8/8/8/8/8/8/8 w - - 0 1');
let lastMove;

let lastMoveBaselineLeft, lastMoveBaselineTop;
let lastMoveBaselineWidth, lastMoveBaselineHeight;
let lastMoveBaselineTransform, lastMoveBaselineTransformOrigin;

let lastMoveArrow1Left, lastMoveArrow1Top;
let lastMoveArrow1Width, lastMoveArrow1Height;
let lastMoveArrow1Transform, lastMoveArrow1TransformOrigin;

let lastMoveArrow2Left, lastMoveArrow2Top;
let lastMoveArrow2Width, lastMoveArrow2Height;
let lastMoveArrow2Transform, lastMoveArrow2TransformOrigin;

let lastMovePointLeft, lastMovePointTop;
let lastMovePointWidth, lastMovePointHeight;
let lastMovePointTransform, lastMovePointTransformOrigin;

let promotionPending = false;
let pendingPromotionMove = undefined;


$: fileIndexes = [true, "true"].includes(reversed) ? [7,6,5,4,3,2,1,0] : [0,1,2,3,4,5,6,7];
$: rankIndexes = [true, "true"].includes(reversed) ? [0,1,2,3,4,5,6,7] : [7,6,5,4,3,2,1,0];

$: cellsSize = size / 9.0;
$: halfCellsSize = cellsSize * 0.5;

$: rootStyle = `
    width: ${size}px;
    height: ${size}px;
    background-color: ${background};
`;

$: gridTemplate = `${halfCellsSize}px repeat(8, ${cellsSize}px) ${halfCellsSize}px / ${halfCellsSize}px repeat(8, ${cellsSize}px) ${halfCellsSize}px`;

$: lowestLayerStyle = `
    width: ${size}px;
    height: ${size}px;
    grid-template: ${gridTemplate};
`;

$: dndLayerStyle = `
    width: ${size}px;
    height: ${size}px;
`;

$: promotionOverlayStyle = `
    width: ${size}px;
    height: ${size}px;
`;

$: promotionStyle = `
    width: ${size * 0.8}px;
    height: ${size * 0.6}px;
`;

$: whiteCellsStyle = `
    background-color: ${white_cell_color};
`;

$: blackCellsStyle = `
    background-color: ${black_cell_color};
`;

$: coordinateStyle = `
    font-size: ${cellsSize * 0.3}px;
    font-weight: bold;
    color: ${coordinates_color};
`;

$: dndOriginCellStyle = `
    background-color: ${origin_cell_color};
    width: ${cellsSize}px;
    height: ${cellsSize}px;
`;

$: dndTargetCellStyle = `
    background-color: ${target_cell_color};
    width: ${cellsSize}px;
    height: ${cellsSize}px;
`;

$: dndCrossCellStyle = `
    background-color: ${dnd_cross_color};
    width: ${cellsSize}px;
    height: ${cellsSize}px;
`;

$: lastMoveBaseLineStyle = `
    left: ${lastMoveBaselineLeft || 0}px;
    top: ${lastMoveBaselineTop || 0}px;
    width: ${lastMoveBaselineWidth || 0}px;
    height: ${lastMoveBaselineHeight || 0}px;
    transform: ${lastMoveBaselineTransform || ''};
    -ms-transform: ${lastMoveBaselineTransform || ''};
    -moz-transform: ${lastMoveBaselineTransform || ''};
    -webkit-transform: ${lastMoveBaselineTransform || ''};
    transform-origin: ${lastMoveBaselineTransformOrigin || '0px 0px'}; 
    -ms-transform-origin: ${lastMoveBaselineTransformOrigin || '0px 0px'};
    -moz-transform-origin: ${lastMoveBaselineTransformOrigin || '0px 0px'}; 
    -webkit-transform-origin: ${lastMoveBaselineTransformOrigin || '0px 0px'};
    background-color: ${move_highlight_color};
`;

$: lastMoveArrow1Style = `
    left: ${lastMoveArrow1Left || 0}px;
    top: ${lastMoveArrow1Top || 0}px;
    width: ${lastMoveArrow1Width || 0}px;
    height: ${lastMoveArrow1Height || 0}px;
    transform: ${lastMoveArrow1Transform || ''};
    -ms-transform: ${lastMoveArrow1Transform || ''};
    -moz-transform: ${lastMoveArrow1Transform || ''};
    -webkit-transform: ${lastMoveArrow1Transform || ''};
    transform-origin: ${lastMoveArrow1TransformOrigin || '0px 0px'}; 
    -ms-transform-origin: ${lastMoveArrow1TransformOrigin || '0px 0px'};
    -moz-transform-origin: ${lastMoveArrow1TransformOrigin || '0px 0px'}; 
    -webkit-transform-origin: ${lastMoveArrow1TransformOrigin || '0px 0px'};
    background-color: ${move_highlight_color};
`;

$: lastMoveArrow2Style = `
    left: ${lastMoveArrow2Left || 0}px;
    top: ${lastMoveArrow2Top || 0}px;
    width: ${lastMoveArrow2Width || 0}px;
    height: ${lastMoveArrow2Height || 0}px;
    transform: ${lastMoveArrow2Transform || ''};
    -ms-transform: ${lastMoveArrow2Transform || ''};
    -moz-transform: ${lastMoveArrow2Transform || ''};
    -webkit-transform: ${lastMoveArrow2Transform || ''};
    transform-origin: ${lastMoveArrow2TransformOrigin || '0px 0px'}; 
    -ms-transform-origin: ${lastMoveArrow2TransformOrigin || '0px 0px'};
    -moz-transform-origin: ${lastMoveArrow2TransformOrigin || '0px 0px'}; 
    -webkit-transform-origin: ${lastMoveArrow2TransformOrigin || '0px 0px'};
    background-color: ${move_highlight_color};
`;

$: lastMovePointStyle = `
    left: ${lastMovePointLeft || 0}px;
    top: ${lastMovePointTop || 0}px;
    width: ${lastMovePointWidth || 0}px;
    height: ${lastMovePointHeight || 0}px;
    transform: ${lastMovePointTransform || ''};
    -ms-transform: ${lastMovePointTransform || ''};
    -moz-transform: ${lastMovePointTransform || ''};
    -webkit-transform: ${lastMovePointTransform || ''};
    transform-origin: ${lastMovePointTransformOrigin || '0px 0px'}; 
    -ms-transform-origin: ${lastMovePointTransformOrigin || '0px 0px'};
    -moz-transform-origin: ${lastMovePointTransformOrigin || '0px 0px'}; 
    -webkit-transform-origin: ${lastMovePointTransformOrigin || '0px 0px'};
    background-color: ${move_highlight_color};
`;

$: halfThickness = cellsSize * 0.08;

$: promotionTitleStyle = `
    margin: ${cellsSize * 0.4}px;
    font-size: ${cellsSize * 0.5}px;
`;

$: promotionButtonStyle = `
    border: ${cellsSize * 0.08}px solid black;
    padding: ${cellsSize * 0.01}px;
    margin: ${cellsSize * 0.2}px;
`;

let dragAndDropInProgress;
let dndPieceData;
let targetFile, targetRank;
let dndLocation;

$: dndPieceStyle = [null, undefined].includes(dndLocation) ? '' : `
    position: absolute;
    left: ${dndLocation.x}px;
    top: ${dndLocation.y}px;
`;

function cancelDnd() {
    dragAndDropInProgress = false;
    dndPieceData = undefined;
    dndLocation = undefined;
    targetFile = undefined;
    targetRank = undefined;
}

function updateWaitingForExternalMove() {
    if (!gameInProgress) return;
    const whiteTurn = logic.turn() === 'w';
    waitingForExternalMove = (whiteTurn && ![true, "true"].includes(white_player_human)) ||
        (!whiteTurn && ![true, "true"].includes(black_player_human));

    if (waitingForExternalMove) dispatch('waiting_manual_move');
}

export function newGame(position = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1') {
    logic = new Chess(position);
    lastMove = undefined;
    promotionPending = false;
    cancelDnd();
    updateWaitingForExternalMove();
    gameInProgress = true;
}

export function toggleSide() {
    reversed = !reversed;
}

export function isWhiteTurn() {
    return logic.turn() === 'w';
}

export function getCurrentPosition() {
    return logic.fen();
}

export function playMove({ startFile, startRank, endFile, endRank, promotion = 'q'}) {
    if (!gameInProgress) return;
    if (!waitingForExternalMove) return;

    const moveObject = {
        from: cellAlgebraic(startFile, startRank),
        to: cellAlgebraic(endFile, endRank),
        promotion,
    }

    const logicBeforeMove = new Chess(logic.fen());

    const result = logic.move(moveObject);
    // Illegal move
    if (! result) return;

    // Update the logic variable => update the board !
    logic = logic;

    updateAndEmitLastMove({
        startFile, startRank,
        endFile, endRank,
        logicBeforeMove,
        logicAfterMove: logic,
    });
    
    handleGameEndedStatus();
    updateWaitingForExternalMove();
}

export function setPositionAndLastMove({
    positionFen,
    fromFileIndex, fromRankIndex,
    toFileIndex, toRankIndex,
}) {
    if (gameInProgress) return;

    logic = new Chess(positionFen);

    lastMove = {
        start: {
            file: fromFileIndex, rank: fromRankIndex,
        },
        end: {
            file: toFileIndex, rank: toRankIndex,
        }
    };
}

function setPromotionPending({startFile, startRank, endFile, endRank}) {
    promotionPending = true;
    pendingPromotionMove = {
        startFile, startRank,
        endFile, endRank,
    };
}

function commitPromotionMove(type) {
    const moveObject = {
        from: cellAlgebraic(pendingPromotionMove.startFile, pendingPromotionMove.startRank),
        to: cellAlgebraic(pendingPromotionMove.endFile, pendingPromotionMove.endRank),
        promotion: type,
    }

    const logicBeforeMove = new Chess(logic.fen());

    logic.move(moveObject);
    // Update the logic variable => update the board !
    logic = logic;

    cancelDnd();
    updateAndEmitLastMove({...pendingPromotionMove,
     logicBeforeMove, logicAfterMove: logic});

    pendingPromotionMove = undefined;
    promotionPending = false;

    handleGameEndedStatus();
    updateWaitingForExternalMove();
}

function updateLogic() {
    // Update the logic variable => update the board !
    logic = logic;
}

function setupDnd({x, y, file, rank, piece}) {
    dndLocation = {
        x, y,
    };

    dndPieceData = {
        ...piece,
        originCell: {file, rank},
    };

    targetFile = file;
    targetRank = rank;

    dragAndDropInProgress = true;
}

function updateDndLocation(x, y, file, rank) {
    dndLocation = {
        x, y,
    };

    targetFile = file;
    targetRank = rank;
}

$: {
    if (![null, undefined].includes(lastMove)) {
        const startColumn = ["true", true].includes(reversed) ? 7 - lastMove.start.file : lastMove.start.file;
        const startLine = ["true", true].includes(reversed) ? lastMove.start.rank : 7 - lastMove.start.rank;
        const endColumn = ["true", true].includes(reversed) ? 7 - lastMove.end.file: lastMove.end.file;
        const endLine = ["true", true].includes(reversed) ? lastMove.end.rank : 7 - lastMove.end.rank;

        const ax = cellsSize * (startColumn + 1.0);
        const ay = cellsSize * (startLine + 1.0);
        const bx = cellsSize * (endColumn + 1.0);
        const by = cellsSize * (endLine + 1.0);

        const realAx = ax - halfThickness;
        const realAy = ay;
        const realBx = bx - halfThickness;
        const realBy = by;

        const vectX = realBx - realAx;
        const vectY = realBy - realAy;

        const baseLineAngleRad = Math.atan2(vectY, vectX) - Math.PI / 2.0;
        const baseLineLength = Math.sqrt(vectX * vectX + vectY * vectY);
        lastMoveBaselineLeft = realAx;
        lastMoveBaselineTop = realAy;
        lastMoveBaselineWidth = 2 * halfThickness;
        lastMoveBaselineHeight = baseLineLength;
        lastMoveBaselineTransform = `rotate(${baseLineAngleRad}rad)`;
        lastMoveBaselineTransformOrigin = `${halfThickness}px ${0}px`;

        const arrow1AngleRad = Math.atan2(vectY, vectX) - Math.PI / 2.0 -  3 * Math.PI / 4.0;
        const arrow1Length = Math.sqrt(vectX * vectX + vectY * vectY) * 0.4;
        lastMoveArrow1Left = realBx;
        lastMoveArrow1Top = realBy;
        lastMoveArrow1Width = 2 * halfThickness;
        lastMoveArrow1Height = arrow1Length;
        lastMoveArrow1Transform = `rotate(${arrow1AngleRad}rad)`;
        lastMoveArrow1TransformOrigin = `${halfThickness}px ${0}px`;

        const arrow2AngleRad = Math.atan2(vectY, vectX) - Math.PI / 2.0 +  3 * Math.PI / 4.0;
        const arrow2Length = Math.sqrt(vectX * vectX + vectY * vectY) * 0.4;
        lastMoveArrow2Left = realBx;
        lastMoveArrow2Top = realBy;
        lastMoveArrow2Width = 2 * halfThickness;
        lastMoveArrow2Height = arrow2Length;
        lastMoveArrow2Transform = `rotate(${arrow2AngleRad}rad)`;
        lastMoveArrow2TransformOrigin = `${halfThickness}px ${0}px`;

        const pointAngleRad = Math.atan2(vectY, vectX) + Math.PI / 4.0;
        const pointLength = 2 * halfThickness;
        lastMovePointLeft = realBx;
        lastMovePointTop = realBy;
        lastMovePointWidth = 2 * halfThickness;
        lastMovePointHeight = pointLength;
        lastMovePointTransform = `rotate(${pointAngleRad}rad)`;
        lastMovePointTransformOrigin = 'center';
    }
    else {
        lastMoveBaselineLeft = undefined; 
        lastMoveBaselineTop = undefined;
        lastMoveBaselineWidth = undefined; 
        lastMoveBaselineHeight = undefined;
        lastMoveBaselineTransform = undefined; 
        lastMoveBaselineTransformOrigin = undefined;

        lastMoveArrow1Left = undefined; 
        lastMoveArrow1Top = undefined;
        lastMoveArrow1Width = undefined; 
        lastMoveArrow1Height = undefined;
        lastMoveArrow1Transform = undefined; 
        lastMoveArrow1TransformOrigin = undefined;

        lastMoveArrow2Left = undefined; 
        lastMoveArrow2Top = undefined;
        lastMoveArrow2Width = undefined; 
        lastMoveArrow2Height = undefined;
        lastMoveArrow2Transform = undefined; 
        lastMoveArrow2TransformOrigin = undefined;

        lastMovePointLeft = undefined; 
        lastMovePointTop = undefined;
        lastMovePointWidth = undefined; 
        lastMovePointHeight = undefined;
        lastMovePointTransform = undefined; 
        lastMovePointTransformOrigin = undefined;
    }
}

function updateAndEmitLastMove({logicBeforeMove, logicAfterMove, startFile, startRank, endFile, endRank}){
    lastMove = {
        start: {
            file: startFile, rank: startRank,
        },
        end: {
            file: endFile, rank: endRank,
        }
    };

    const moveNumber = logicBeforeMove.fen().split(" ")[5];
    const allMovesHistory = logicAfterMove.history();
    const whiteTurnBeforeMove = logicBeforeMove.turn() === 'w';
    const moveSan = allMovesHistory[allMovesHistory.length - 1];
    const moveFan = convertMoveSanToMoveFan({moveSan, whiteTurn: whiteTurnBeforeMove});

    const lastMoveEventPayload = {
        moveNumber,
        whiteTurn: whiteTurnBeforeMove,
        moveSan,
        moveFan,
        positionFen: logicAfterMove.fen(),
        fromFileIndex: startFile,
        fromRankIndex: startRank,
        toFileIndex: endFile,
        toRankIndex: endRank,
    };

    dispatch('move_done', {moveObject: lastMoveEventPayload});
}

function convertMoveSanToMoveFan({moveSan, whiteTurn}) {
    moveSan = moveSan.replace(/K/g, whiteTurn ? '\u2654' : '\u265A').normalize("NFKC");
    moveSan = moveSan.replace(/Q/g, whiteTurn ? '\u2655' : '\u265B').normalize("NFKC");
    moveSan = moveSan.replace(/R/g, whiteTurn ? '\u2656' : '\u265C').normalize("NFKC");
    moveSan = moveSan.replace(/B/g, whiteTurn ? '\u2657' : '\u265D').normalize("NFKC");
    moveSan = moveSan.replace(/N/g, whiteTurn ? '\u2658' : '\u265E').normalize("NFKC");

    return moveSan;
};

function handleGameEndedStatus() {
    if (logic.in_checkmate()) {
        cancelDnd();
        gameInProgress = false;
        dispatch('checkmate', {whiteTurnBeforeMove: logic.turn() !== 'w'});
    }
    else if (logic.in_stalemate()) {
        cancelDnd();
        gameInProgress = false;
        dispatch('stalemate');
    }
    else if (logic.in_threefold_repetition()) {
        cancelDnd();
        gameInProgress = false;
        dispatch('perpetual_draw');
    }
    else if (logic.in_draw()) {
        if (logic.insufficient_material()) {
            cancelDnd();
            gameInProgress = false;
            dispatch('missing_material_draw');
        }
        else {
            cancelDnd();
            gameInProgress = false;
            dispatch('fifty_moves_draw');
        }
    }
}
</script>

<style>
.root {
    position: relative;
}

.lowest-layer {
    position: absolute;
    display: grid;
    z-index: 2;
}

.last-move-layer {
    position: absolute;
    z-index: 3;
}

.dnd-layer {
    position: absolute;
    z-index: 4;
}

.promotion-overlay-layer {
    position: absolute;
    z-index: 5;
    opacity: 0.8;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: black;
}

.promotion-layer {
    position: absolute;
    background-color: white;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
}

.promotion-title {
    font-weight: bold;
    color: black;
}

.promotion-buttons {
    display: flex;
    justify-content: center;
    align-items: center;
}

.cell {
    display: flex;
    justify-content: center;
    align-items: center;
}

.coordinate {
    display: flex;
    justify-content: center;
    align-items: center;
}

.dndPiece {
    position: absolute;
}

.last-move-line {
    position: absolute;
}

.player-turn {
    border-radius: 50%;
}

.white-turn {
    background-color: white;
}

.black-turn {
    background-color: black;
}
</style>

<svelte:options tag="loloof64-chessboard" />
<div class="root" bind:this={rootElement} style={rootStyle}
     on:mousedown|preventDefault={(event) => handleMouseDown({event, cellsSize, reversed, rootElement, 
        logic, dragAndDropInProgress, setupDnd, gameInProgress, waitingForExternalMove})}
     on:mousemove|preventDefault={(event) => handleMouseMove({event, dragAndDropInProgress,
        updateDndLocation, rootElement, cancelDnd, cellsSize, reversed, 
        promotionPending, gameInProgress, waitingForExternalMove})}
     on:mouseup|preventDefault={(event) => handleMouseUp({event, cellsSize, reversed, rootElement,
        logic, dragAndDropInProgress, dndPieceData, cancelDnd, updateLogic,
         updateAndEmitLastMove, promotionPending, setPromotionPending, gameInProgress, 
         handleGameEndedStatus, updateWaitingForExternalMove, waitingForExternalMove})}
     on:mouseleave|preventDefault={(event) => handleMouseExited({event, cancelDnd, promotionPending, 
        gameInProgress, waitingForExternalMove})}
>
    <div class="lowest-layer" style={lowestLayerStyle}>
        <div></div>
        {#each fileIndexes as file}
            <div class="coordinate"  style={coordinateStyle}>{String.fromCharCode('A'.charCodeAt(0) + file)}</div>
        {/each}
        <div></div>

        {#each rankIndexes as rank}
            <div class="coordinate"  style={coordinateStyle}>{String.fromCharCode('1'.charCodeAt(0) + rank)}</div>
            {#each fileIndexes as file}
                <div class="cell" style="{((rank + file) % 2)  === 0 ? whiteCellsStyle : blackCellsStyle}">
                    {#if isDnDOriginCell(dndPieceData, file, rank)}
                        <div 
                            style={
                                targetFile !== file || targetRank !== rank ? dndOriginCellStyle : dndTargetCellStyle
                            }
                                
                        ></div>
                    {:else if isWhitePawnAtCell(logic, file, rank)}
                        <chess-white-pawn size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle  :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isWhiteKnightAtCell(logic, file, rank)}
                        <chess-white-knight size={cellsSize}
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                         />
                    {:else if isWhiteBishopAtCell(logic, file, rank)}
                        <chess-white-bishop size={cellsSize}
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                         />
                    {:else if isWhiteRookAtCell(logic, file, rank)}
                        <chess-white-rook size={cellsSize}
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                         />
                    {:else if isWhiteQueenAtCell(logic, file, rank)}
                        <chess-white-queen size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isWhiteKingAtCell(logic, file, rank)}
                        <chess-white-king size={cellsSize}
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                         />

                    {:else if isBlackPawnAtCell(logic, file, rank)}
                        <chess-black-pawn size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isBlackKnightAtCell(logic, file, rank)}
                        <chess-black-knight size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isBlackBishopAtCell(logic, file, rank)}
                        <chess-black-bishop size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isBlackRookAtCell(logic, file, rank)}
                        <chess-black-rook size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isBlackQueenAtCell(logic, file, rank)}
                        <chess-black-queen size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else if isBlackKingAtCell(logic, file, rank)}
                        <chess-black-king size={cellsSize} 
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        />
                    {:else}
                        <div
                            style={
                                targetFile === file && targetRank === rank ? dndTargetCellStyle :
                                (targetFile === file || targetRank === rank) ? dndCrossCellStyle : ''}
                        ></div>
                    {/if}
                </div>
            {/each}
            <div class="coordinate" style={coordinateStyle}>{String.fromCharCode('1'.charCodeAt(0) + rank)}</div>
        {/each}

        <div></div>
        {#each fileIndexes as file}
            <div class="coordinate" style={coordinateStyle}>{String.fromCharCode('A'.charCodeAt(0) + file)}</div>
        {/each}
        <div class="player-turn {logic.turn() === 'w' ? 'white-turn' : 'black-turn'}"></div>
    </div>

    <div class="dnd-layer">
        {#if isWhitePawnDragged(dndPieceData)}
            <chess-white-pawn size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isWhiteKnightDragged(dndPieceData)}
            <chess-white-knight size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isWhiteBishopDragged(dndPieceData)}
            <chess-white-bishop size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isWhiteRookDragged(dndPieceData)}
            <chess-white-rook size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isWhiteQueenDragged(dndPieceData)}
            <chess-white-queen size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isWhiteKingDragged(dndPieceData)}
            <chess-white-king size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isBlackPawnDragged(dndPieceData)}
            <chess-black-pawn size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isBlackKnightDragged(dndPieceData)}
            <chess-black-knight size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isBlackBishopDragged(dndPieceData)}
            <chess-black-bishop size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isBlackRookDragged(dndPieceData)}
            <chess-black-rook size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isBlackQueenDragged(dndPieceData)}
            <chess-black-queen size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else if isBlackKingDragged(dndPieceData)}
            <chess-black-king size={cellsSize} style={dndPieceStyle} class="dndPiece"/>
        {:else}
            <div></div>
        {/if}
    </div>

    <div class="last-move-layer" style={dndLayerStyle}>
        <div class="last-move-line" style={lastMoveBaseLineStyle}></div>
        <div class="last-move-line" style={lastMoveArrow1Style}></div>
        <div class="last-move-line" style={lastMoveArrow2Style}></div>
        <div class="last-move-line" style={lastMovePointStyle}></div>
    </div>

    {#if promotionPending === true}
        <div class="promotion-overlay-layer" style={promotionOverlayStyle}>
            <div class="promotion-layer" style={promotionStyle}>
                <div class="promotion-title" style={promotionTitleStyle}>{promotion_dialog_title}</div>
                <div class="promotion-buttons">
                    {#if logic.turn() === 'w'}
                        <chess-white-queen 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('q')}
                        />
                        <chess-white-rook 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('r')}
                        />
                        <chess-white-bishop 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('b')} 
                        />
                        <chess-white-knight 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('n')} 
                        />
                    {:else}
                        <chess-black-queen 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('q')}
                        />
                        <chess-black-rook 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('r')} 
                        />
                        <chess-black-bishop 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('b')} 
                        />
                        <chess-black-knight 
                            style={promotionButtonStyle} size={cellsSize}
                            on:click={() => commitPromotionMove('n')}
                        />
                    {/if}
                </div>
            </div>
        </div>
    {/if}
</div>