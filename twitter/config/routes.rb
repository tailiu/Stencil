Rails.application.routes.draw do
	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
	root 'pages#show'

	# get 'pages/login', to: 'pages#login'
	# get 'pages/signUp', to: 'pages#signUp'
	# get 'pages/home', to: 'pages#home'
	# get 'pages/messages', to: 'pages#messages'
	# get 'pages/notifications', to: 'pages#notifs'
	# get 'pages/search', to: 'pages#search'
	# get 'pages/profile', to: 'pages#profile'
	# get 'pages/settings', to: 'pages#settings'

	get 'users/new', to: 'users#new'
	get 'users/verify', to: 'users#verify'
	get 'users/logout', to: 'users#logout'

	resources :users

	resources :pages
end
