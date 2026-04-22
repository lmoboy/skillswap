import { writable } from 'svelte/store'

export interface User {
   name: string
   email: string
   id: number | string
   profile_picture?: string
   is_admin?: boolean
}

interface AuthState {
   user: User | null
   isAuthenticated: boolean
   loading: boolean
   error: string | null
   step: string | null
}

const defaultState: AuthState = {
   user: null,
   isAuthenticated: false,
   loading: true,
   error: null,
   step: null,
}

// Helper function to persist auth state to localStorage
function persistAuthState(state: AuthState) {
   try {
      if (typeof window !== 'undefined') {
         const persistableState = {
            user: state.user,
            isAuthenticated: state.isAuthenticated,
            // Don't persist loading, error, step as they're transient
         }
         localStorage.setItem('authState', JSON.stringify(persistableState))
      }
   } catch (e) {
      // Ignore storage errors
      console.warn('Failed to persist auth state:', e)
   }
}

// Helper function to restore auth state from localStorage
function restoreAuthState(): AuthState {
   try {
      if (typeof window !== 'undefined') {
         const stored = localStorage.getItem('authState')
         if (stored) {
            const parsed = JSON.parse(stored)
            return {
               ...defaultState,
               ...parsed,
               loading: false, // Don't restore loading state
            }
         }
      }
   } catch (e) {
      // Ignore restoration errors
      console.warn('Failed to restore auth state:', e)
   }
   return defaultState
}

function createAuthStore() {
   // Initialize with persisted state or default
   const initialState = restoreAuthState()
   const { subscribe, set, update } = writable<AuthState>(initialState)

   // Subscribe to state changes for persistence
   const unsubscribe = subscribe((state) => {
      persistAuthState(state)
   })

   // Cleanup subscription on unload
   if (typeof window !== 'undefined') {
      window.addEventListener('beforeunload', () => {
         unsubscribe()
      })
   }

   return {
      subscribe,
      setUser: (user: User | null) => {
         update((state) => {
            const newState = {
               ...state,
               user,
               isAuthenticated: user !== null,
               loading: false,
               error: null,
            }
            return newState
         })
      },
      setStep: (step: string) => {
         update((state) => ({ ...state, step }))
      },
      setLoading: (loading: boolean) => {
         update((state) => ({ ...state, loading }))
      },
      setError: (error: string | null) => {
         update((state) => ({ ...state, error, loading: false }))
      },
      clearUser: () => {
         const newState = { ...defaultState, loading: false }
         set(newState)
      },
      isAuthenticated: () => {
         let isAuth = false
         const unsubscribe = subscribe((state) => {
            isAuth = state.isAuthenticated
         })
         unsubscribe()
         return isAuth
      },
      /** Wait until user is available or timeout */
      waitForUser: (timeoutMs = 5000): Promise<User> => {
         return new Promise((resolve, reject) => {
            let unsub: () => void
            const timer = setTimeout(() => {
               unsub?.()
               reject(new Error('Timed out waiting for user'))
            }, timeoutMs)

            unsub = subscribe((state) => {
               if (!state.loading && state.user) {
                  clearTimeout(timer)
                  unsub()
                  resolve(state.user)
               }
            })
         })
      },
      // Method to reset the entire store to default state
      reset: () => {
         set(defaultState)
      },
   }
}

export const auth = createAuthStore()
export const { setUser, clearUser, reset } = auth
