Rails.application.routes.draw do
	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html

	get 'users/verify', to: 'users#verify'
	get 'users/logout', to: 'users#logout'
	get 'users/getFollowers', to: 'users#getFollowers'

	resources :users

end
