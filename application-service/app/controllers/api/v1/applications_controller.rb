module Api
  module V1
    class ApplicationsController < ApplicationController
      def create
        application = Application.new(application_params)

        if application.save
          render json: {
            token: application.token,
            name: application.name
          }, status: :created
        else
          render json: { errors: application.errors }, status: :unprocessable_entity
        end
      end

      def index
        applications = Application.all
        render json: applications, status: :ok
      end
      
      def show
        token = params[:id]
        Rails.logger.debug "Received token: #{token}"
        begin
          application = Application.find_by(token: token)
          render json: application, status: :ok
        rescue ActiveRecord::RecordNotFound
          render json: { error: 'Application not found' }, status: :not_found
        end
      end
      private

      def application_params
        params.require(:application).permit(:name)
      end
    end
  end
end
