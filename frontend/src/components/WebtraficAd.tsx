import React, { useEffect } from 'react';

const WebtraficAd: React.FC = () => {
    const adContainerRef = React.useRef<HTMLDivElement>(null);

    useEffect(() => {
        const container = adContainerRef.current;
        if (!container) return;

        // Check if script is already in this container to prevent duplicates
        if (container.querySelector('script[src="https://webtrafic.ru/ads.php?uid=17708"]')) {
            return;
        }

        const script = document.createElement('script');
        script.src = "https://webtrafic.ru/ads.php?uid=17708";
        script.async = true;

        container.appendChild(script);

        return () => {
            // Optional: cleanup
        };
    }, []);

    return (
        <div
            id="webtraf_17708"
            ref={adContainerRef}
            style={{
                width: '468px',
                height: '60px',
                display: 'inline-block', // Ensure it behaves well in flex
                backgroundColor: 'rgba(255,255,255,0.1)', // Temporary placeholder background to see if div exists
            }}
            className="webtrafic-ad"
        />
    );
};

export default WebtraficAd;
