Rails.application.routes.draw do
	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
	get 'pages/login', to: 'pages#login'
	get 'pages/signUp', to: 'pages#signUp'

	resources :users

	resources :pages
end
