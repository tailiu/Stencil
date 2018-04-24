class UsersController < ApplicationController
    def new
        @new_user = User.create(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.credential.create()

        if @new_user.save
            render "pages/home"
        else 
            puts @new_user.errors.messages
            render "pages/err"
        end

        # @existing_user = User.find_by name: params

        # if @new_user
        #     render "pages/home"
        # else
        #     render "pages/login"
        # end

    end
end
