<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTheme } from '../composables/useTheme'
import { useOidcAuth } from '#imports'
import { label } from '@primeuix/themes/aura/metergroup'


const { is_light, toggleTheme } = useTheme()
const { user, isLoggedIn, login, logout } = useOidcAuth()

const isCollapsed = ref(false)
const isMobileMenuOpen = ref(false)

const userNameInitial = computed(() => {
  if (user.value && user.value.profile && user.value.profile.name) {
    return user.value.profile.name.charAt(0).toUpperCase()
  }
  return ''
})

const navItems = computed(() => {
  if (isLoggedIn()) {
    return menuItems
  }

  return menuItems
    .map(group => ({
      ...group,
      items: group.items.filter(item => item.protected === false),
    }))
    .filter(group => group.items.length > 0)
})


const menuItems = [
  {
    label: 'Home',
    group: "",
    items: [
      {
        label: 'Home',
        id: "/",
        icon: 'bi bi-house',
        protected: false
      }
    ]
  },
  {
    label: 'IoT Connectors',
    group: "connectors",
    items: [
      {
        label: 'Add',
        icon: 'bi bi-plus-lg',
        protected: true
      },
      {
        label: 'Manage',
        icon: 'bi bi-wrench',
        protected: true
      },
      {
        label: 'Documentation',
        icon: 'bi bi-book',
        protected: true
      }
    ]
  },
  {
    label: 'Assets',
    group: "assets",
    items: [
      {
        label: 'Add',
        icon: 'bi bi-plus-lg',
        protected: true
      },
      {
        label: 'Manage',
        icon: 'bi bi-wrench',
        protected: true
      },
      {
        label: 'Documentation',
        icon: 'bi bi-book',
        protected: true
      }
    ]
  },
  {
    label: 'Dashboards',
    group: "dashboards",
    items: [
      {
        label: 'Design',
        id: "design",
        icon: 'bi bi-brush',
        protected: true
      },
      {
        label: 'View',
        icon: 'bi bi-eye',
        protected: true
      }
    ]
  },
  {
    label: 'Project',
    group: "project",
    items: [
      {
        label: 'Settings',
        icon: 'bi bi-gear',
        protected: true
      }
    ]
  }
];


onMounted(() => {
  menuItems.forEach(g => {
    g.items.forEach(i => {
      i.route = `/${g.group}/${i.id}`
    })
  })

})




function toggleMobileMenu() {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}


function routing(path: string, props: any) {
  console.log(path, props)
}

</script>

<template>
  <div>
    <button class="sidebar-toggle" @click="toggleMobileMenu">
      <i class="bi bi-list" />
    </button>
    <aside class="sidebar" :class="{ collapsed: isCollapsed, 'mobile-open': isMobileMenuOpen }">
      <nav class="sidebar-nav">

        <div class="logo-container">
          <span class="logo-small">G</span>
        </div>
        <div class="login-container" >
          <Avatar :label="userNameInitial" size="large" v-if="isLoggedIn()"/>
          <Button v-if="!isLoggedIn()" icon="bi bi-box-arrow-in-right" :fluid="true" label="Login" @click="login()" />
          <Button v-else icon="bi bi-box-arrow-in-left" :fluid="true" label="Logout" @click="logout()"/>
        </div>


        <Menu :model="navItems" id="nav-menu">
          <template #item="{ item, props }">

            <a v-bind="props.action" class="menu-item" @click="routing(item.route, props)">
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
.p-avatar {
  border: 3px solid var(--a-color-1);
  box-shadow: 0px 0px 10px var(--a-color-1)
}


.logo-container {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  font-family: "Zen Dots";
  border-bottom: 2px solid var(--p-background);
  padding: 5px;
}

.logo-container>span {
  margin-right: 0px;
  font-size: 2rem;
}


.login-container {
  display: flex;
  padding: 1rem;
  flex-direction: column;

  align-items: center;
  justify-content: center;
}

.login-container>button {
  margin-top: 4px;
  background: var(--p-background);
  border: none;
  color: var(--p-text-color);
  padding: 5px 10px;
}

.login-container>button:hover {
  background-color: var(--p-background) !important;
  border: none !important;
  color: var(--p-text-color) !important;
}

.logo-small {
  font-size: 2rem;
  display: inline-block;
  font-family: "Zen Dots";
  transform: rotate(45deg);
  text-shadow: -1px -2px 3px rgb(255, 0, 119), -1px -5px 3px rgb(0, 140, 255);

}

.menu-item {
  border: 1px solid var(--p-color-2);
  border-radius: 4px;
  margin: 1px;
  background: var(--p-background);
}

.menu-item:hover {
  border-left: 3px solid var(--a-color-1);
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