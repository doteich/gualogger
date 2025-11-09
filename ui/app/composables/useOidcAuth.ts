// composables/useOidcAuth.ts
import { UserManager, User } from 'oidc-client-ts';

// We create a single instance of the UserManager
let userManager: UserManager;

const getManager = () => {
    let config = useRuntimeConfig()


    if (!userManager) {
        const settings = {
            authority: config.public.authority, // Your Keycloak realm URL
            client_id: config.public.clientId, // The Client ID you set in Keycloak
            redirect_uri: config.public.redirectUrl, // The callback URL
            post_logout_redirect_uri: config.public.redirectUrl,
            response_type: 'code', // Use Authorization Code Flow
            scope: config.public.scopes, // Standard OIDC scopes
            loadUserInfo: true, // Optional: auto-loads user profile after login
        };
        userManager = new UserManager(settings);
    }
    return userManager;
};

export const useOidcAuth = () => {
    const manager = getManager();
    const user = useState<User | null>('oidc-user', () => null);

    // Check for a user on initial load
    manager.getUser().then(result => {
        user.value = result;
    });

    const login = () => {
        return manager.signinRedirect();
    };

    const logout = () => {
        return manager.signoutRedirect();
    };

    const isLoggedIn = () => {
        return user.value?.access_token ? true : false

    }

    // Expose the raw manager to handle callbacks
    return {
        manager,
        user,
        login,
        logout,
        isLoggedIn
    };
};