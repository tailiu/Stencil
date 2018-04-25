class UsersController < ApplicationController
    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            render "pages/home"
        else 
            puts @new_user.errors.messages
            puts @new_credential.errors.messages
            render "pages/err"

        end

    end

end
