
/**
 * Parse synopsis content that might be in JSON format from upstream scrapers.
 * Example format: {"connections": [], "paragraphs": ["Para 1", "Para 2"]}
 */
export function parseSynopsis(content: string | any): string {
  if (!content) return "No description available.";
  
  // If it's already an object (not stringified)
  if (typeof content === "object") {
    if (content.paragraphs && Array.isArray(content.paragraphs)) {
        return content.paragraphs.join("\n\n");
    }
    return JSON.stringify(content); // Fallback
  }

  // If it's a string, try to detect JSON
  if (typeof content === "string") {
    const trimmed = content.trim();
    if (trimmed.startsWith("{") && trimmed.endsWith("}")) {
      try {
        const parsed = JSON.parse(content);
        if (parsed.paragraphs && Array.isArray(parsed.paragraphs)) {
          return parsed.paragraphs.join("\n\n");
        }
      } catch (e) {
        // Not valid JSON, return original string
      }
    }
  }

  return content;
}
