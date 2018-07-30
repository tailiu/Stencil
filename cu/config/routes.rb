Rails.application.routes.draw do
  get 'query/index'
  get 'query/submit', to: 'query#submit'
  get 'query/result', to: 'query#result'
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
end
