Rails.application.routes.draw do
	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html

	get 'users/verify', to: 'users#verify'
	get 'users/logout', to: 'users#logout'
	get 'tweets/fetchall', to: 'tweets#fetchall'
	get 'tweets/fetchUserTweets', to: 'tweets#fetchUserTweets'
	get 'tweets/mainPageTweets', to: 'tweets#mainPageTweets'
	get 'users/getUserInfo', to: 'users#getUserInfo'

	resources :users
	resources :user_actions
	resources :tweets
end
