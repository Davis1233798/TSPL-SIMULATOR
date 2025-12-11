import React, { useEffect } from 'react';

const WebtraficAd: React.FC = () => {
    const iframeRef = React.useRef<HTMLIFrameElement>(null);

    useEffect(() => {
        const iframe = iframeRef.current;
        if (!iframe) return;

        const doc = iframe.contentDocument || iframe.contentWindow?.document;
        if (!doc) return;

        // Ensure we write to an empty document
        doc.open();
        doc.write(`
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <style>
                    body { margin: 0; padding: 0; overflow: hidden; background: transparent; }
                    /* Center content if needed, though ad is usually fixed size */
                </style>
            </head>
            <body>
                <div id="webtraf_17708"></div>
                <!-- Synchronous script to ensure it executes immediately in this fresh context -->
                <script src="https://webtrafic.ru/ads.php?uid=17708"></script>
            </body>
            </html>
        `);
        doc.close();
    }, []);

    return (
        <iframe
            ref={iframeRef}
            title="Sponsor Ad"
            style={{
                width: '468px',
                height: '60px',
                border: 'none',
                overflow: 'hidden',
                backgroundColor: 'transparent',
                display: 'inline-block'
            }}
            className="webtrafic-ad"
            scrolling="no"
        />
    );
};

export default WebtraficAd;
