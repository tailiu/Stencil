class UsersController < ApplicationController
    def new
        @user = User.new(params)
        puts params[:a]

        render "pages/home"
    end
end
