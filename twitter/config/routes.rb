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
	get 'users/search', to: 'users#search'
	get 'users/updateBio', to: 'users#updateBio'
	get 'users/checkFollow', to: 'user_actions#checkFollow'
	get 'users/checkMute', to: 'user_actions#checkMute'
	get 'users/checkBlock', to: 'user_actions#checkBlock'
	get 'users/checkTwoWayBlock', to: 'user_actions#checkTwoWayBlock'
	get 'users/handleMute', to: 'user_actions#handleMute'
	get 'users/handleBlock', to: 'user_actions#handleBlock'
	get 'users/handleFollow', to: 'user_actions#handleFollow'
	get 'users/approveFollowRequest', to: 'user_actions#approveFollowRequest'
	get 'users/getFollowRequests', to: 'user_actions#getFollowRequests'
	get 'users/updateEmail', to: 'users#updateEmail'
	get 'users/updateHandle', to: 'users#updateHandle'
	get 'users/updatePassword', to: 'users#updatePassword'
	get 'users/updateProtected', to: 'users#updateProtected'
	post 'tweets/newf', to: 'tweets#new'
	get 'tweets/delete', to: 'tweets#delete'
	get 'notifications/get', to: 'notifications#get'

	get 'conversations/leaveConversation', to: 'conversations#leaveConversation'
	get 'conversations/blockInGroupConversation', to: 'conversations#blockInGroupConversation'

	resources :users
	resources :user_actions
	resources :tweets
	resources :conversations
	resources :messages


	
end
