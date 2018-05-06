Rails.application.routes.draw do
	# For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html

	get 'users/verify', to: 'users#verify'
	get 'users/logout', to: 'users#logout'
	get 'tweets/fetchall', to: 'tweets#fetchall'
	get 'tweets/fetchUserTweets', to: 'tweets#fetchUserTweets'
	get 'tweets/mainPageTweets', to: 'tweets#mainPageTweets'
	get 'tweets/getTweet', to: 'tweets#getTweet'
	get 'tweets/like', to: 'tweets#like'
	get 'tweets/retweet', to: 'tweets#retweet'
	get 'tweets/stats', to: 'tweets#stats'
	get 'users/getUserInfo', to: 'users#getUserInfo'
	get 'users/updateBio', to: 'users#updateBio'
	get 'users/checkFollow', to: 'user_actions#checkFollow'
	get 'users/handleFollow', to: 'user_actions#handleFollow'
	get 'users/updateEmail', to: 'users#updateEmail'
	get 'users/updateHandle', to: 'users#updateHandle'
	get 'users/updatePassword', to: 'users#updatePassword'
	get 'users/updateProtected', to: 'users#updateProtected'

	resources :users
	resources :user_actions
	resources :tweets
	resources :conversations
	resources :messages


	
end
