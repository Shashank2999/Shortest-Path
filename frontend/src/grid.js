import React, { useState, useEffect } from 'react';
import './grid.css';

const Grid = () => {
    const [selectedCells, setSelectedCells] = useState({ start: null, end: null });
    const [path, setPath] = useState([]);

    useEffect(() => {
        if (selectedCells.start && selectedCells.end) {
            fetchPath(selectedCells.start, selectedCells.end);
        }
    }, [selectedCells]);

    const handleCellClick = (x, y) => {
        if (!selectedCells.start) {
            setSelectedCells({ ...selectedCells, start: { x, y } });
        } else if (!selectedCells.end) {
            setSelectedCells({ ...selectedCells, end: { x, y } });
        }
    };

    const fetchPath = async (start, end) => {
        try {
            const res = await fetch('http://localhost:8080/find-path', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ start, end }),
            });
            const data = await res.json();
            setPath(data.shortestPath);
        } catch (error) {
            console.error('Error fetching path:', error);
        }
    };

    const handleReset = () => {
        setSelectedCells({ start: null, end: null });
        setPath([]);
    };

    return (
        <div className="grid-container">
            <div className="grid">
                {[...Array(20)].map((_, row) => (
                    <div key={row} className="row">
                        {[...Array(20)].map((_, col) => {
                            const isStart = selectedCells.start?.x === row && selectedCells.start?.y === col;
                            const isEnd = selectedCells.end?.x === row && selectedCells.end?.y === col;
                            const isPath = path.some(p => p.x === row && p.y === col);
                            return (
                                <div
                                    key={col}
                                    className={`cell ${isStart ? 'start' : isEnd ? 'end' : isPath ? 'path' : ''}`}
                                    onClick={() => handleCellClick(row, col)}
                                >
                                </div>
                            );
                        })}
                    </div>
                ))}
            </div>
            <button className="reset-button" onClick={handleReset}>
                Reset
            </button>
        </div>
    );
};

export default Grid;