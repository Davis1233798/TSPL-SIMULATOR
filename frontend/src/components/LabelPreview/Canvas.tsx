import React, { useEffect, useRef } from 'react';
import { RenderData, RenderElement } from '../../types/tspl';

interface CanvasProps {
  renderData: RenderData;
}

const Canvas: React.FC<CanvasProps> = ({ renderData }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  // 常數定義
  const DPI = 203; // 點/英寸
  const MM_TO_INCH = 0.0393701; // mm 轉英寸
  const SCALE = 2; // 縮放比例以提高清晰度

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // 計算畫布尺寸（以點為單位）
    const widthInDots = Math.round(renderData.width * MM_TO_INCH * DPI);
    const heightInDots = Math.round(renderData.height * MM_TO_INCH * DPI);

    // 設定畫布尺寸
    canvas.width = widthInDots * SCALE;
    canvas.height = heightInDots * SCALE;

    // 縮放上下文
    ctx.scale(SCALE, SCALE);

    // 繪製白色背景
    ctx.fillStyle = 'white';
    ctx.fillRect(0, 0, widthInDots, heightInDots);

    // 繪製邊框
    ctx.strokeStyle = '#ccc';
    ctx.lineWidth = 1;
    ctx.strokeRect(0, 0, widthInDots, heightInDots);

    // 繪製元素
    renderData.elements.forEach((element) => {
      drawElement(ctx, element);
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [renderData]);

  const drawElement = (ctx: CanvasRenderingContext2D, element: RenderElement) => {
    ctx.save();

    switch (element.type) {
      case 'TEXT':
        drawText(ctx, element);
        break;
      case 'BARCODE':
        drawBarcode(ctx, element);
        break;
      case 'QRCODE':
        drawQRCode(ctx, element);
        break;
      case 'BOX':
        drawBox(ctx, element);
        break;
      case 'BAR':
        drawBar(ctx, element);
        break;
    }

    ctx.restore();
  };

  const drawText = (ctx: CanvasRenderingContext2D, element: RenderElement) => {
    const { x, y, properties } = element;
    const { content, y_scale = 1, rotation = 0 } = properties;

    ctx.fillStyle = 'black';
    ctx.font = `${16 * y_scale}px monospace`;

    ctx.translate(x, y);
    if (rotation) {
      ctx.rotate((rotation * Math.PI) / 180);
    }

    ctx.fillText(content, 0, 0);
  };

  const drawBarcode = (ctx: CanvasRenderingContext2D, element: RenderElement) => {
    const { x, y, properties } = element;
    const { content, height = 100, narrow = 2 } = properties;

    ctx.fillStyle = 'black';

    // 簡化的條碼繪製（實際應該根據條碼類型生成）
    const barWidth = narrow;
    let currentX = x;

    for (let i = 0; i < content.length; i++) {
      if (i % 2 === 0) {
        ctx.fillRect(currentX, y, barWidth, height);
      }
      currentX += barWidth * 2;
    }

    // 繪製條碼文字
    ctx.font = '12px monospace';
    ctx.textAlign = 'center';
    ctx.fillText(content, x + (currentX - x) / 2, y + height + 15);
  };

  const drawQRCode = (ctx: CanvasRenderingContext2D, element: RenderElement) => {
    const { x, y, properties } = element;
    const { cell_width = 5 } = properties;

    // 簡化的 QR Code 繪製（實際應該使用 QR Code 生成庫）
    const size = cell_width * 25; // 假設 25x25 模塊

    ctx.fillStyle = 'white';
    ctx.fillRect(x, y, size, size);

    ctx.fillStyle = 'black';
    ctx.strokeStyle = 'black';
    ctx.lineWidth = 1;
    ctx.strokeRect(x, y, size, size);

    // 繪製簡化的 QR 碼圖案
    for (let i = 0; i < 25; i++) {
      for (let j = 0; j < 25; j++) {
        if (Math.random() > 0.5) {
          ctx.fillRect(
            x + i * cell_width,
            y + j * cell_width,
            cell_width,
            cell_width
          );
        }
      }
    }
  };

  const drawBox = (ctx: CanvasRenderingContext2D, element: RenderElement) => {
    const { x, y, properties } = element;
    const { width, height, thickness = 1 } = properties;

    ctx.strokeStyle = 'black';
    ctx.lineWidth = thickness;
    ctx.strokeRect(x, y, width, height);
  };

  const drawBar = (ctx: CanvasRenderingContext2D, element: RenderElement) => {
    const { x, y, properties } = element;
    const { width, height } = properties;

    ctx.fillStyle = 'black';
    ctx.fillRect(x, y, width, height);
  };

  return (
    <div className="canvas-container">
      <canvas ref={canvasRef} className="label-canvas" />
    </div>
  );
};

export default Canvas;
