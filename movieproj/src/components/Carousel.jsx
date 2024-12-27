import React, { useState, useEffect } from "react";
import './Caraousel.css';
import image1 from './voPyk0.webp';
import image2 from './329583.jpg'
import image3 from './1561771.jpg'
import image4 from './wp9164084.jpg'

const Carousel = () => {
  const [currentIndex, setCurrentIndex] = useState(0);


  const images = [
    `${image1}`, `${image2}`, `${image3}`, `${image4}`
    ,


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
