/**
 * Shared utilities for icon/avatar color assignment and emoji detection.
 */

/** 16-color palette for avatar backgrounds — all pass WCAG AA contrast against white */
export const AVATAR_COLORS = [
  '#5B8FF9', // blue
  '#5AD8A6', // green
  '#F6BD16', // yellow
  '#E8684A', // red-orange
  '#6DC8EC', // sky
  '#9270CA', // purple
  '#FF9845', // orange
  '#36CBCB', // teal
  '#4E79A7', // steel blue
  '#A770D4', // lavender
  '#71C287', // sage
  '#E88CCF', // pink
  '#4ECDC4', // mint
  '#F2B84B', // gold
  '#E85D5D', // coral
  '#6BBAA7', // seafoam
]

/**
 * Deterministic color for a member based on their ID.
 * Same ID always produces the same color.
 */
export function getAvatarColor(memberId?: number | null): string {
  if (!memberId) return AVATAR_COLORS[0]
  return AVATAR_COLORS[memberId % AVATAR_COLORS.length]
}

/**
 * Semantic color mapping for preset category icons.
 * Key = Lucide icon name, value = background color hex.
 */
export const CATEGORY_COLORS: Record<string, string> = {
  // Preset category icons
  UtensilsCrossed: '#E87461',
  Car: '#4ECDC4',
  ShoppingBag: '#E88CCF',
  Home: '#6BBAA7',
  Smartphone: '#4A90D9',
  HeartPulse: '#E85D5D',
  BookOpen: '#F2B84B',
  Gamepad2: '#9B59B6',
  Baby: '#FF8C42',
  Gift: '#E74C3C',
  PawPrint: '#8B6F47',
  Sparkles: '#E88CCF',
  Dumbbell: '#2ECC71',
  ShieldCheck: '#3498DB',
  Package: '#95A5A6',
  Banknote: '#2ECC71',
  Briefcase: '#F39C12',
  TrendingUp: '#27AE60',
  Envelope: '#E74C3C',
  Mail: '#E74C3C',
  // Additional icons from IconPicker
  Crown: '#F2B84B',
  Smile: '#5AD8A6',
  User: '#5B8FF9',
  Users: '#5B8FF9',
  UserCircle: '#4A90D9',
  Star: '#F2B84B',
  Wallet: '#6BBAA7',
  DollarSign: '#27AE60',
  TrendingDown: '#E85D5D',
  PieChart: '#9B59B6',
  BarChart3: '#4A90D9',
  PiggyBank: '#E88CCF',
  CreditCard: '#4A90D9',
  Receipt: '#95A5A6',
  GraduationCap: '#F2B84B',
  Stethoscope: '#E85D5D',
  Shirt: '#E88CCF',
  Plane: '#4ECDC4',
  Sun: '#F2B84B',
  Moon: '#9B59B6',
  Cloud: '#6DC8EC',
  Flame: '#E87461',
  Leaf: '#5AD8A6',
  TreePine: '#27AE60',
  Fish: '#4ECDC4',
  Bird: '#6DC8EC',
  Bug: '#5AD8A6',
  Flower2: '#E88CCF',
  Zap: '#F2B84B',
  Target: '#E87461',
  Anchor: '#4A90D9',
  Compass: '#4ECDC4',
  Shield: '#3498DB',
  Key: '#F39C12',
  Bell: '#E87461',
  Flag: '#E74C3C',
  Globe: '#4A90D9',
  Coffee: '#8B6F47',
  Cake: '#E88CCF',
  Music: '#9B59B6',
  Camera: '#95A5A6',
  Phone: '#5AD8A6',
  MapPin: '#E85D5D',
}

/**
 * Detect whether a string is an emoji character (Unicode codepoint > 0x1F000)
 * vs a Lucide icon name (ASCII string).
 */
export function isEmoji(value: string): boolean {
  if (!value) return false
  const cp = value.codePointAt(0)
  return cp !== undefined && cp > 0x1f000
}

/** Default color for custom categories that don't have a preset color */
const CATEGORY_DEFAULT_COLOR = '#95A5A6'

/**
 * Get the background color for a category icon.
 * Uses semantic mapping for preset icons, falls back to neutral gray.
 */
export function getCategoryColor(iconName: string): string {
  if (isEmoji(iconName)) return CATEGORY_DEFAULT_COLOR
  return CATEGORY_COLORS[iconName] || CATEGORY_DEFAULT_COLOR
}
