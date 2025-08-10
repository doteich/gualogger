<script setup lang="ts">
import { ref } from 'vue'
import { useTheme } from '../composables/useTheme'

const { is_light, toggleTheme } = useTheme()

const isCollapsed = ref(false)
const isMobileMenuOpen = ref(false)

const menuItems = ref([
    {
    label: 'Home',
    items: [
      {
        label: 'Home',
        icon: 'bi bi-house'
      }
    ]
  },
  {
    label: 'IoT Connectors',
    items: [
      {
        label: 'Add',
        icon: 'bi bi-plus-lg'
      },
      {
        label: 'Manage',
        icon: 'bi bi-wrench'
      },
      {
        label: 'Documentation',
        icon: 'bi bi-book'
      }
    ]
  },
  {
    label: 'Assets',
    items: [
      {
        label: 'Add',
        icon: 'bi bi-plus-lg'
      },
      {
        label: 'Manage',
        icon: 'bi bi-wrench'
      },
      {
        label: 'Documentation',
        icon: 'bi bi-book'
      }
    ]
  },
  {
    label: 'Dashboards',
    items: [
      {
        label: 'Design',
        icon: 'bi bi-brush'
      },
      {
        label: 'View',
        icon: 'bi bi-eye'
      }
    ]
  },
  {
    label: 'Project',
    items: [
      {
        label: 'Settings',
        icon: 'bi bi-gear'
      }
    ]
  }
]);

function toggleMobileMenu() {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}
</script>

<template>
  <div>
    <button class="sidebar-toggle" @click="toggleMobileMenu">
      <i class="bi bi-list" />
    </button>
    <aside class="sidebar" :class="{ collapsed: isCollapsed, 'mobile-open': isMobileMenuOpen }">
      <nav class="sidebar-nav">
        <Menu :model="menuItems" id="nav-menu">
          <template #item="{ item, props }">
            <a v-bind="props.action" class="menu-item">
              <span :class="item.icon" class="menu-icon" />
              <span>{{ item.label }}</span>

            </a>
          </template>
        </Menu>
        <div class="theme-toggle-container">
          <ToggleButton v-model="is_light" offIcon="bi bi-sun-fill" onIcon="bi bi-moon-fill" onLabel="Dark Mode"
            offLabel="Light Mode" fluid />

        </div>
      </nav>
    </aside>
  </div>
</template>

<style>
.menu-item {
  border: 1px solid var(--p-color-2);
  border-radius: 4px;
  margin: 1px;
  background: var(--p-background);
}

.menu-icon {
  background: var(--p-color-2);
  padding: 3px 5px;
  border-radius: 4px;

}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.theme-toggle-container {
  margin-top: auto;
  display: flex;

  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.sidebar {
  background: var(--p-color-1);
  height: 100vh;
  padding: 5px;
}

#nav-menu {
  border: none;
  width: inherit;
  background: none;

}

#nav-menu>* {
  width: inherit !important;
}


.theme-toggle {
  display: flex;
  align-items: center;
  width: 100%;
  background: none;
  border: none;
  color: var(--p-text-color);
  padding: 1rem;
  cursor: pointer;
  font-size: 1rem;
  border-radius: 4px;
}

.theme-toggle:hover {
  background-color: var(--p-color-2);
}

.theme-toggle .p-menuitem-label {
  margin-left: 0.5rem;
}

.sidebar.collapsed .theme-toggle .p-menuitem-label {
  display: none;
}

.sidebar.collapsed .theme-toggle {
  justify-content: center;
}

.sidebar {
  background: var(--p-color-1);
  transition: width 0.3s ease;
}

#nav-menu {
  border: none;
  width: 100%;
}

#nav-menu .p-menuitem {
  width: 100%;
}

#nav-menu .p-menuitem-link {
  display: flex;
  align-items: center;
  padding: 1rem;
  color: var(--p-text-color);
  text-decoration: none;
  transition: background-color 0.3s;
}

#nav-menu .p-menuitem-link:hover {
  background-color: var(--p-color-2);
}

#nav-menu .p-menuitem-icon {
  margin-right: 0.5rem;
}

.sidebar-toggle {
  display: none;
  position: fixed;
  top: 1rem;
  left: 1rem;
  z-index: 1001;
  background: var(--p-color-1);
  border: none;
  color: var(--p-text-color);
  padding: 0.5rem;
  border-radius: 4px;
  cursor: pointer;
}

/* Mobile Styles */
@media (max-width: 768px) {
  .sidebar-toggle {
    display: block;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    height: 100%;
    width: 250px;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    z-index: 1000;
  }

  .sidebar.mobile-open {
    transform: translateX(0);
  }

  .sidebar.collapsed {
    width: 250px;
    /* Override collapsed state on mobile */
  }
}

/* Desktop Styles */
@media (min-width: 769px) {
  .sidebar {
    width: 250px;
  }

  .sidebar.collapsed {
    width: 60px;
  }

  .sidebar.collapsed .p-menuitem-label {
    display: none;
  }

  .sidebar.collapsed .p-menuitem-link {
    justify-content: center;
  }

  .sidebar.collapsed .p-menuitem-icon {
    margin-right: 0;
  }
}
</style>