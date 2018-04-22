Rails.application.routes.draw do
	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
	get 'loginOrSignup', to: 'pages#loginOrSignUp'
	
	resources :users

	resources :pages
end
