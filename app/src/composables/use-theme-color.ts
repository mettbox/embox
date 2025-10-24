export const useThemeColor = () => {
  /**
   * Set the theme color for the current route.
   *
   * @param {string} color - The color to set (can be a CSS variable or a direct color value).
   */
  const setThemeColor = (color: string) => {
    // Check if the color is a CSS variable
    if (color.startsWith('--')) {
      const rootStyles = getComputedStyle(document.documentElement);

      color = rootStyles.getPropertyValue(color).trim();
    }

    const metaTag = document.querySelector('meta[name="theme-color"]');

    if (metaTag) {
      metaTag.setAttribute('content', color);
    } else {
      const newMetaTag = document.createElement('meta');
      newMetaTag.setAttribute('name', 'theme-color');
      newMetaTag.setAttribute('content', color);
      document.head.appendChild(newMetaTag);
    }
  };

  return { setThemeColor };
};
