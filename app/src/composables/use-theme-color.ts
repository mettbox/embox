export const useThemeColor = () => {
  /**
   * Set the theme color for the current route.
   *
   * @param {string} color - The color to set (can be a CSS variable or a direct color value).
   */
  const setThemeColor = (color: string) => {
    if (color.startsWith('--')) {
      const rootStyles = getComputedStyle(document.documentElement);
      color = rootStyles.getPropertyValue(color).trim();
    }

    document.body.style.backgroundColor = color;
  };

  return { setThemeColor };
};
