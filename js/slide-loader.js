/**
 * Slide Loader System for MCP Presentation
 * Dynamically loads individual slide files and assembles them into the presentation
 */

class SlideLoader {
  constructor(manifestPath = 'slides/slides-manifest.json', slidesContainer = '.slides') {
    this.manifestPath = manifestPath;
    this.slidesContainer = slidesContainer;
    this.slides = [];
    this.manifest = null;
  }

  /**
   * Load the slide manifest
   */
  async loadManifest() {
    try {
      const response = await fetch(this.manifestPath);
      if (!response.ok) {
        throw new Error(`Failed to load manifest: ${response.status}`);
      }
      this.manifest = await response.json();
      console.log(`Loaded manifest: ${this.manifest.title}`);
      return this.manifest;
    } catch (error) {
      console.error('Error loading manifest:', error);
      throw error;
    }
  }

  /**
   * Load a single slide file
   */
  async loadSlide(slideConfig) {
    try {
      const response = await fetch(`slides/${slideConfig.file}`);
      if (!response.ok) {
        throw new Error(`Failed to load slide ${slideConfig.file}: ${response.status}`);
      }
      const html = await response.text();
      return {
        ...slideConfig,
        content: html
      };
    } catch (error) {
      console.error(`Error loading slide ${slideConfig.file}:`, error);
      // Return a placeholder slide for missing files
      return {
        ...slideConfig,
        content: `<section><h2>Error loading slide: ${slideConfig.title}</h2><p>Could not load ${slideConfig.file}</p></section>`
      };
    }
  }

  /**
   * Load all slides concurrently
   */
  async loadAllSlides() {
    if (!this.manifest) {
      await this.loadManifest();
    }

    console.log(`Loading ${this.manifest.slides.length} slides...`);
    
    // Load all slides concurrently for better performance
    const slidePromises = this.manifest.slides.map(slideConfig => 
      this.loadSlide(slideConfig)
    );

    try {
      this.slides = await Promise.all(slidePromises);
      console.log(`Successfully loaded ${this.slides.length} slides`);
      return this.slides;
    } catch (error) {
      console.error('Error loading slides:', error);
      throw error;
    }
  }

  /**
   * Inject slides into the DOM
   */
  injectSlides() {
    const container = document.querySelector(this.slidesContainer);
    if (!container) {
      throw new Error(`Slides container '${this.slidesContainer}' not found`);
    }

    // Clear existing content
    container.innerHTML = '';

    // Add each slide to the container
    this.slides.forEach((slide, index) => {
      const slideElement = document.createElement('div');
      slideElement.innerHTML = slide.content;
      
      // Extract the section element (since each slide file contains a full section)
      const sectionElement = slideElement.querySelector('section');
      if (sectionElement) {
        // Add data attributes for debugging and navigation
        sectionElement.setAttribute('data-slide-id', slide.id);
        sectionElement.setAttribute('data-slide-index', index);
        sectionElement.setAttribute('data-slide-title', slide.title);
        
        container.appendChild(sectionElement);
      } else {
        console.warn(`No section element found in slide ${slide.file}`);
        // Fallback: add the content as-is
        container.appendChild(slideElement);
      }
    });

    console.log(`Injected ${this.slides.length} slides into the DOM`);
  }

  /**
   * Initialize the slide loader and load all slides
   */
  async initialize() {
    try {
      console.log('Initializing slide loader...');
      await this.loadAllSlides();
      this.injectSlides();
      console.log('Slide loader initialization complete');
      return true;
    } catch (error) {
      console.error('Failed to initialize slide loader:', error);
      return false;
    }
  }

  /**
   * Get slide information
   */
  getSlideInfo(slideId) {
    return this.slides.find(slide => slide.id === slideId);
  }

  /**
   * Get all slide metadata
   */
  getManifest() {
    return this.manifest;
  }
}

// Export for use in other scripts
window.SlideLoader = SlideLoader;