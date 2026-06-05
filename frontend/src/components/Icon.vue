<script setup lang="ts">
import { computed } from 'vue'
import {
  LayoutDashboard,
  Wallet,
  ListTodo,
  Star,
  MessageSquare,
  Users,
  FolderOpen,
  UserCircle,
  Settings,
  Moon,
  Sun,
  TrendingUp,
  TrendingDown,
  Plus,
  Trash2,
  Pencil,
  Search,
  Filter,
  Calendar,
  DollarSign,
  Tag,
  MessageCircle,
  Heart,
  ArrowLeft,
  MoreHorizontal,
  X,
  Check,
  AlertTriangle,
  XCircle,
  Info,
  ChevronLeft,
  ChevronRight,
  Eye,
  EyeOff,
  Lock,
  User,
  Bell,
  Home,
  Send,
  ThumbsUp,
  Clock,
  FileText,
  PieChart,
  BarChart3,
  Minus,
} from 'lucide-vue-next'

const iconMap: Record<string, any> = {
  LayoutDashboard,
  Wallet,
  ListTodo,
  Star,
  MessageSquare,
  Users,
  FolderOpen,
  UserCircle,
  Settings,
  Moon,
  Sun,
  TrendingUp,
  TrendingDown,
  Plus,
  Trash2,
  Pencil,
  Search,
  Filter,
  Calendar,
  DollarSign,
  Tag,
  MessageCircle,
  Heart,
  ArrowLeft,
  MoreHorizontal,
  X,
  Check,
  AlertTriangle,
  XCircle,
  Info,
  ChevronLeft,
  ChevronRight,
  Eye,
  EyeOff,
  Lock,
  User,
  Bell,
  Home,
  Send,
  ThumbsUp,
  Clock,
  FileText,
  PieChart,
  BarChart3,
  Minus,
}

const props = withDefaults(defineProps<{
  name: string
  size?: number
  color?: string
  strokeWidth?: number
}>(), {
  size: 20,
  color: 'currentColor',
  strokeWidth: 1.5,
})

const iconComponent = computed(() => iconMap[props.name])
</script>

<template>
  <component
    v-if="iconComponent"
    :is="iconComponent"
    :size="size"
    :color="color"
    :stroke-width="strokeWidth"
  />
</template>
