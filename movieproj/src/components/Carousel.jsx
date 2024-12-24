import React, { useState, useEffect } from "react";
import './Caraousel.css'; 

const Carousel = () => {
  const [currentIndex, setCurrentIndex] = useState(0);

  
  const images = [
    "https://via.placeholder.com/600x300?text=Image+1",
    "https://via.placeholder.com/600x300?text=Image+2",
    "https://via.placeholder.com/600x300?text=Image+3",
    "https://via.placeholder.com/600x300?text=Image+4",
  ];


  const nextImage = () => {
    setCurrentIndex((prevIndex) => (prevIndex + 1) % images.length);
  };

 
  const prevImage = () => {
    setCurrentIndex((prevIndex) => (prevIndex - 1 + images.length) % images.length);
  };


  useEffect(() => {
    const interval = setInterval(nextImage, 5000);
    return () => clearInterval(interval); 
  }, []);

  return (
    <div className="carousel-container">
      <button className="carousel-button prev" onClick={prevImage}>
        &#10094;
      </button>
      <div className="carousel">
        <img src={images[currentIndex]} alt={`Slide ${currentIndex}`} />
      </div>
      <button className="carousel-button next" onClick={nextImage}>
        &#10095;
      </button>
    </div>
  );
};

export default Carousel;
