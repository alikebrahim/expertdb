/**
 * Basic fuzzy search implementation
 * Calculates similarity between search term and target string
 */
export const fuzzyMatch = (searchTerm: string, target: string): number => {
  if (!searchTerm || !target) return 0;
  
  const search = searchTerm.toLowerCase();
  const text = target.toLowerCase();
  
  // Direct match gets highest score
  if (text.includes(search)) {
    return 1.0;
  }
  
  // Character-by-character matching
  let score = 0;
  let searchIndex = 0;
  
  for (let i = 0; i < text.length && searchIndex < search.length; i++) {
    if (text[i] === search[searchIndex]) {
      score += 1;
      searchIndex++;
    }
  }
  
  // Normalize score by search term length
  return searchIndex === search.length ? score / search.length : 0;
};

/**
 * Fuzzy search function that returns filtered and scored results
 */
export const fuzzySearch = <T>(
  items: T[],
  searchTerm: string,
  getSearchableFields: (item: T) => string[],
  minScore: number = 0.3
): T[] => {
  if (!searchTerm.trim()) return items;
  
  const results = items
    .map(item => {
      const fields = getSearchableFields(item);
      const scores = fields.map(field => fuzzyMatch(searchTerm, field));
      const maxScore = Math.max(...scores, 0);
      
      return { item, score: maxScore };
    })
    .filter(result => result.score >= minScore)
    .sort((a, b) => b.score - a.score)
    .map(result => result.item);
  
  return results;
};

/**
 * Highlight matching characters in text
 */
export const highlightMatch = (text: string, searchTerm: string): string => {
  if (!searchTerm || !text) return text;
  
  const search = searchTerm.toLowerCase();
  const target = text.toLowerCase();
  
  // Simple highlighting for substring matches
  const index = target.indexOf(search);
  if (index !== -1) {
    return (
      text.substring(0, index) +
      '<mark>' +
      text.substring(index, index + searchTerm.length) +
      '</mark>' +
      text.substring(index + searchTerm.length)
    );
  }
  
  return text;
};